package service

import (
	"atium/pkg/atium"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type dbStorage struct {
	db *sql.DB
}

func newStore(username, password, host, dbname string) (*dbStorage, error) {
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbname)
	db, err := sql.Open("mysql", connectString)
	if err == nil {
		return &dbStorage{db}, nil
	}
	return nil, err
}

func (ms *dbStorage) Close() error {
	return ms.db.Close()
}

// -----------------------------------------------------------------------------
// Client related DB function which implements interface 'storage'
// -----------------------------------------------------------------------------
// TODO: listClients
func (ms *dbStorage) listClients() (interface{}, error) {
	return nil, nil
}

func (ms *dbStorage) getClient(name string) (*atium.ClientDetails, error) {
	qs := fmt.Sprintf("%s %s %s",
		"SELECT CL.id, CL.description, CL.created_at, CL.modified_at,",
		"C.email, C.address, C.primary_ph, C.secondary_ph FROM client CL",
		"LEFT JOIN contact C ON C.id = CL.contact_id WHERE name = ?")
	var address string
	var id int64
	c := atium.ClientDetails{}
	c.Name = name
	row := ms.db.QueryRow(qs, name)
	err := row.Scan(&id, &c.Description, &c.CreatedAt,
		&c.ModifiedAt, &c.Email, &address, &c.PrimaryPh, &c.SecondaryPh)
	if err != nil {
		return nil, fmt.Errorf("client not found: %v", err)
	}
	err = json.Unmarshal([]byte(address), &c.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(address)
	}
	services, err := getClientServiceList(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting services: %v", err)
	}
	c.Services = services
	return &c, nil
}

func (ms *dbStorage) getClientInfo(name string) (*atium.ClientInfo, error) {
	c := atium.ClientInfo{}
	row := ms.db.QueryRow("SELECT description, modified_at, created_at FROM client WHERE name = ?", name)
	c.Name = name
	err := row.Scan(&c.Description, &c.ModifiedAt, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("err getting user: %v", err)
	}
	return &c, nil
}

func (ms *dbStorage) getClientUsers(client string) ([]*atium.UserDetails, error) {
	// get client id
	id, err := getClientID(ms.db, client)
	if err != nil {
		return nil, fmt.Errorf("err getting client id: %v", err)
	}
	// get list of users mapped to client
	var result []*atium.UserDetails
	qs := fmt.Sprintf("%s %s",
		"SELECT U.Email FROM user U LEFT JOIN client_user CU ON CU.user_id = U.id",
		"WHERE client_id = ?")
	rows, err := ms.db.Query(qs, id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		user, err := ms.getUser(email)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, nil
}

func (ms *dbStorage) getClientServices(name string) ([]atium.ClientService, error) {
	id, err := getClientID(ms.db, name)
	if err != nil {
		return nil, fmt.Errorf("err getting client id: %v", err)
	}
	return getClientServiceList(ms.db, id)
}

func (ms *dbStorage) getClientActivities(name string) ([]atium.ActivityInfo, error) {
	id, err := getClientID(ms.db, name)
	if err != nil {
		return nil, fmt.Errorf("err getting client id: %v", err)
	}
	return getClientActivititesList(ms.db, id)
}

func (ms *dbStorage) upsertClient(c atium.ClientDetails) error {
	contact := atium.ContactInfo{Email: c.Email, Address: c.Address,
		PrimaryPh: c.PrimaryPh, SecondaryPh: c.SecondaryPh}
	// update contact table
	err := ms.upsertUserContact(c.Email, contact)
	if err != nil {
		return err
	}
	// update client table
	ci := atium.ClientInfo{Description: c.Description, Name: c.Name}
	err = ms.upsertClientInfo(ci)
	if err != nil {
		return err
	}
	return nil
}

func (ms *dbStorage) upsertClientInfo(c atium.ClientInfo) error {
	// update user table
	qs1 := []string{"name = ?"}
	qs2 := []string{}
	args1 := []interface{}{c.Name}
	args2 := []interface{}{}
	if c.Description != "" {
		qs2 = append(qs2, "description = ?")
		args2 = append(args2, c.Description)
	}
	qs1 = append(qs1, qs2...)
	args1 = append(args1, args2...)
	params := append(args1, time.Now())
	params = append(params, args2...)
	qs := fmt.Sprintf("%s %s %s %s", "INSERT user SET",
		strings.Join(qs1, ", "), "ON DUPLICATE KEY UPDATE modified_at = ?", strings.Join(qs2, ", "))
	_, err := ms.db.Exec(qs, params...)
	return err
}

// TODO: deleteClient
func (ms *dbStorage) deleteClient(a string) error {
	return nil
}

// -----------------------------------------------------------------------------
// User related DB function which implements interface 'storage'
// -----------------------------------------------------------------------------

func (ms *dbStorage) getUser(email string) (*atium.UserDetails, error) {
	qs := fmt.Sprintf("%s %s %s %s",
		"SELECT U.id, U.name, U.dob, U.gender, U.created_at,",
		"C.address, C.primary_ph, C.secondary_ph FROM user U",
		"LEFT JOIN contact C ON C.email = U.email",
		"WHERE U.email = ?")

	var address string
	var id int64
	u := atium.UserDetails{}
	row := ms.db.QueryRow(qs, email)
	u.Email = email
	err := row.Scan(&id, &u.Name, &u.Dob, &u.Gender, &u.CreatedAt,
		&address, &u.PrimaryPh, &u.SecondaryPh)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	err = json.Unmarshal([]byte(address), &u.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(address)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

func (ms *dbStorage) getUserInfo(email string) (*atium.UserInfo, error) {
	u := atium.UserInfo{}
	row := ms.db.QueryRow("SELECT name, dob, gender, created_at FROM user WHERE email = ?", email)
	u.Email = email
	err := row.Scan(&u.Name, &u.Dob, &u.Gender, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("err getting user: %v", err)
	}
	return &u, nil
}

func (ms *dbStorage) getUserContact(email string) (*atium.ContactInfo, error) {
	var address string
	qs := "SELECT address, primary_ph, secondary_ph FROM contact WHERE email = ?"
	row := ms.db.QueryRow(qs, email)
	c := atium.ContactInfo{}
	c.Email = email
	err := row.Scan(&address, &c.PrimaryPh, &c.SecondaryPh)
	if err != nil {
		return nil, fmt.Errorf("contact not found: %v", err)
	}
	err = json.Unmarshal([]byte(address), &c.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(c.Address)
	}
	return &c, err
}

func (ms *dbStorage) getUserStats(email string) (*atium.StatsInfo, error) {
	id, err := getUserID(ms.db, email)
	if err != nil {
		return nil, fmt.Errorf("err getting user id: %v", err)
	}
	return getLatestStats(ms.db, id)
}

func (ms *dbStorage) getUserActivities(email, detail string) ([]atium.ActivityInfo, error) {
	id, err := getUserID(ms.db, email)
	if err != nil {
		return nil, fmt.Errorf("err getting client id: %v", err)
	}
	return getUserActivititesList(ms.db, id, detail)
}

func (ms *dbStorage) upsertUser(u atium.UserDetails) error {
	c := atium.ContactInfo{Email: u.Email, Address: u.Address,
		PrimaryPh: u.PrimaryPh, SecondaryPh: u.SecondaryPh}
	// update contact table
	err := ms.upsertUserContact(u.Email, c)
	if err != nil {
		return err
	}

	ui := atium.UserInfo{Email: u.Email, Name: u.Name,
		Dob: u.Dob, Gender: u.Gender}
	// update user table
	err = ms.upsertUserInfo(ui)
	if err != nil {
		return err
	}

	// update stats table
	if u.Stats != (atium.StatsInfo{}) {
		return ms.upsertUserStats(u.Email, u.Stats)
	}
	return nil
}

// TODO: upsertUserActivities
func (ms *dbStorage) upsertUserActivities(a atium.ActivityInfo) error {
	return nil
}

func (ms *dbStorage) upsertUserInfo(u atium.UserInfo) error {
	// update user table
	qs1 := []string{"email = ?"}
	qs2 := []string{}
	args1 := []interface{}{u.Email}
	args2 := []interface{}{}
	if u.Name != "" {
		qs2 = append(qs2, "name = ?")
		args2 = append(args2, u.Name)
	}
	if !u.Dob.IsZero() {
		qs2 = append(qs2, "dob = ?")
		args2 = append(args2, u.Dob)
	}
	if u.Gender != "" {
		qs2 = append(qs2, "gender = ?")
		args2 = append(args2, u.Gender)
	}
	qs1 = append(qs1, qs2...)
	args1 = append(args1, args2...)
	params := append(args1, args2...)
	qs := fmt.Sprintf("%s %s %s %s", "INSERT user SET",
		strings.Join(qs1, ", "), "ON DUPLICATE KEY UPDATE", strings.Join(qs2, ", "))
	_, err := ms.db.Exec(qs, params...)
	return err
}

func (ms *dbStorage) upsertUserStats(email string, s atium.StatsInfo) error {
	userID, err := getUserID(ms.db, email)
	if err != nil {
		return err
	}
	var insertID int64
	qs := []string{"user_id = ?"}
	args := []interface{}{userID}
	if s.Height != 0 {
		qs = append(qs, "height = ?")
		args = append(args, s.Height)
	}
	if s.Weight != 0 {
		qs = append(qs, "weight = ?")
		args = append(args, s.Weight)
	}
	if s.Arms != 0 {
		qs = append(qs, "arms = ?")
		args = append(args, s.Arms)
	}
	if s.Chest != 0 {
		qs = append(qs, "chest = ?")
		args = append(args, s.Chest)
	}
	if s.Waist != 0 {
		qs = append(qs, "waist = ?")
		args = append(args, s.Waist)
	}
	if s.Hips != 0 {
		qs = append(qs, "hips = ?")
		args = append(args, s.Hips)
	}
	if s.Thighs != 0 {
		qs = append(qs, "thighs = ?")
		args = append(args, s.Thighs)
	}
	if s.Calves != 0 {
		qs = append(qs, "calves = ?")
		args = append(args, s.Calves)
	}
	if len(qs) > 1 {
		qs1 := fmt.Sprintf("%s %s", "INSERT stats SET ",
			strings.Join(qs, ", "))
		r, err := ms.db.Exec(qs1, args...)
		if err != nil {
			return err
		}
		insertID, err = r.LastInsertId()
		if err != nil {
			return err
		}
	}
	if insertID != 0 {
		_, err = ms.db.Exec("INSERT INTO user_stats (user_id, stats_id) VALUES (?, ?)",
			userID, insertID)
		if err != nil {
			return err
		}
	}
	return err
}

func (ms *dbStorage) upsertUserContact(email string, c atium.ContactInfo) error {
	if c.Email == "" {
		c.Email = email
	} else if email != c.Email {
		return fmt.Errorf("updating email is not supported")
	}
	// update contact table
	qs1 := []string{"email = ?"}
	qs2 := []string{}
	args1 := []interface{}{c.Email}
	args2 := []interface{}{}
	var address string
	b, err := json.Marshal(c.Address)
	if err != nil {
		fmt.Println("marshalling address failed", c.Address)
	} else {
		address = string(b)
	}
	if address != "" {
		qs2 = append(qs2, "address = ?")
		args2 = append(args2, address)
	}
	if c.PrimaryPh != "" {
		qs2 = append(qs2, "primary_ph = ?")
		args2 = append(args2, c.PrimaryPh)
	}
	if c.SecondaryPh != "" {
		qs2 = append(qs2, "secondary_ph = ?")
		args2 = append(args2, c.SecondaryPh)
	}
	qs1 = append(qs1, qs2...)
	args1 = append(args1, args2...)
	params := append(args1, args2...)
	qs := fmt.Sprintf("%s %s %s %s", "INSERT contact SET",
		strings.Join(qs1, ", "), "ON DUPLICATE KEY UPDATE", strings.Join(qs2, ", "))
	_, err = ms.db.Exec(qs, params...)
	return err
}

// TODO: deleteUser
func (ms *dbStorage) deleteUser(a string) error {
	return nil
}

// -----------------------------------------------------------------------------
// Service related DB functions which implements interface 'storage'
// -----------------------------------------------------------------------------

func (ms *dbStorage) listServices() ([]atium.ServiceInfo, error) {
	var services []atium.ServiceInfo
	rows, err := ms.db.Query("SELECT S.name, S.description, S.price FROM service")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		s := atium.ServiceInfo{}
		err = rows.Scan(&s.Name, &s.Description, &s.Price)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func (ms *dbStorage) getService(name string) (*atium.ServiceInfo, error) {
	row := ms.db.QueryRow(
		"SELECT S.name, S.description, S.price FROM service WHERE name = ?", name)
	s := atium.ServiceInfo{}
	err := row.Scan(&s.Name, &s.Description, &s.Price)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// -----------------------------------------------------------------------------
// Activity related DB functions which implements interface 'storage'
// -----------------------------------------------------------------------------

func (ms *dbStorage) listActivities() ([]atium.ActivityInfo, error) {
	var activities []atium.ActivityInfo
	rows, err := ms.db.Query(
		"SELECT name, description, type, duration, cost, default_forms FROM activity")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		formsIDStr := ""
		a := atium.ActivityInfo{}
		err = rows.Scan(&a.Name, &a.Description, &a.Type, &a.Duration, &a.Cost, &formsIDStr)
		if err != nil {
			return nil, err
		}
		var formIds []int64
		for _, id := range strings.Split(formsIDStr, ",") {
			fid, _ := strconv.ParseInt(id, 10, 64)
			formIds = append(formIds, fid)
		}
		a.Forms, err = getFormsDetails(ms.db, formIds)
		if err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}

func (ms *dbStorage) getActivity(name string) (*atium.ActivityInfo, error) {
	var a atium.ActivityInfo
	formsIDStr := ""
	row := ms.db.QueryRow(
		"SELECT description, type, duration, cost, default_forms FROM activity WHERE name = ?", name)
	err := row.Scan(&a.Description, &a.Type, &a.Duration, &a.Cost, &formsIDStr)
	if err != nil {
		return nil, fmt.Errorf("activity not found: %v", err)
	}
	a.Name = name
	var formIds []int64
	for _, id := range strings.Split(formsIDStr, ",") {
		fid, _ := strconv.ParseInt(id, 10, 64)
		formIds = append(formIds, fid)
	}
	a.Forms, err = getFormsDetails(ms.db, formIds)
	if err != nil {
		return nil, err
	}
	return &a, err
}

// -----------------------------------------------------------------------------
// Form related DB functions which implements interface 'storage'
// -----------------------------------------------------------------------------

func (ms *dbStorage) listForms() ([]atium.Form, error) {
	var forms []atium.Form
	rows, err := ms.db.Query("SELECT name, label, type, template FROM form")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		f := atium.Form{}
		err = rows.Scan(&f.Name, &f.Label, &f.Type, &f.Template)
		if err != nil {
			return nil, err
		}
		forms = append(forms, f)
	}
	return forms, nil
}

func (ms *dbStorage) getForm(label string) (*atium.Form, error) {
	var f atium.Form
	row := ms.db.QueryRow(
		"SELECT name, type, template FROM form WHERE label = ?", label)
	err := row.Scan(&f.Name, &f.Type, &f.Template)
	f.Label = label
	if err != nil {
		return nil, fmt.Errorf("form not found: %v", err)
	}
	return &f, err
}

// -----------------------------------------------------------------------------
// Consumption related DB functions which implements interface 'storage'
// -----------------------------------------------------------------------------

// TODO: uncomment after finalising FoodInfo struct

// func (ms *dbStorage) getConsumption(foodID string) (*atium.FoodInfo, error) {
// 	qs := fmt.Sprintf("%s %s %s %s",
// 		"SELECT calories_in, protein, carbs, fat, entry_at FROM consumption C",
// 		"WHERE food_id = ?")
// 	var id int64
// 	f := atium.FoodInfo{}
// 	f.FoodID = foodID
// 	row := ms.db.QueryRow(qs, email)
// 	err := row.Scan(&f.CaloriesIn, &f.Protein, &f.Carbs, &f.Fat, &f.EntryAt)
// 	if err != nil {
// 		return nil, fmt.Errorf("consumption not found: %v", err)
// 	}
// 	return &f, nil
// }
