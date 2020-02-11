package service
import (
  "atium/pkg/atium"
  "database/sql"
  "encoding/json"
  "fmt"
  "_github.com/go-sql-driver/mysql"
  "string"
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


  func handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  	dbHost := os.Getenv("RDS_HOST")
  	dbUser := os.Getenv("RDS_USER")
  	dbPass := os.Getenv("RDS_PASSWORD")
  	dbName := os.Getenv("RDS_DB_NAME")
  	// AWS Credentials not required to im
  	var err error
  	db, err := newStore(dbUser, dbPass, dbHost, dbName)
  	if err == nil {
  		c.DB = db
  	}
  	defer c.DB.Close()

// 1.User
func (ms *dbStorage) getUser(name string) (*atium.UserInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "SELECT U.id, U.name, C.email, U.dob, U.gender, U.created_at, C.address,",
			  "C.primary_ph, C.secondary_ph from user U",
			  "LEFT JOIN contact C ON C.id = U.contact_id",
			  "WHERE name = ?")
	var address string
	var id int64
	u := atium.UserInfo{}
	row := ms.db.QueryRow(qs, email)
	err := row.Scan(&id, &u.Name, &u.Email, &u.Dob, &u.Gender, &u.CreatedAt,
			&address, &u.PrimaryPh, &u.SecondaryPh)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	err = json.Unmarshal([]byte(address), &u.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(u.Address)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

// 2.Client
func (ms *dbStorage) getclient(email string) (*atium.clientInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select L.id,L.name,C.email, L.services,L.description,L.created_at,L.modified_at,C.address,",
			  "C.primary_ph,C.secondary_ph from client L",
			  "LEft JOIN contact C ON C.id= L.contact_id",
			  "where email = ?")
	var address string
	var services string
	var description string
	var id int64
	l := atium.clientInfo{}
	row := ms.db.QueryRow(qs, email)
	err := row.Scan(&id,&l.Name,&l.email,&services,&description,&l.created_at,&l.modified_at,&address,
			&l.primary_ph,&l.secondary_ph)
	if err != nil {
		return nil, fmt.Errorf("client not found: %v", err)
	}
	err = json.Unmarshal([]byte(Address), &l.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(l.Address)
	}
	err = json.Unmarshal([]byte( services), &l.services)
	if err != nil {
		fmt.Printf("unmarshalling  services failed: %v", err)
		fmt.Println(l.services)
	}
	err = json.Unmarshal([]byte(description), &l.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(l.description)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	l.Stats = *stats
	return &l, err
}

//3.service
func (ms *dbStorage) getservice(name string) (*atium.serviceInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select S.id,C.client.id,S.name,S.description,S.price,",
			  "C.enabledAt,C.enabled from service S",
			  "LEFT JOIN client_service C ON C.service_id = S.id",
			  "WHERE name = ?")
	var description string
	var id int64
	u := atium.serviceInfo{}
	row := ms.db.QueryRow(qs, email)
	err := row.Scan(&id,&u.client_id,&u.name,&description,&u.price,
			&u.enabledAt,&u.enabled )
	if err != nil {
		return nil, fmt.Errorf("service not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &u.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(u.description)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//4.client_service
func (ms *dbStorage) getclient_service(name string) (*atium.client_serviceInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
	"select C.id,C.name,C.description,C.services,C.created_at,C.modified_at,",
			  "S.enabledAt,S.enabled from client_service S",
			  "LEft JOIN client C ON C.id = S.client_id",
			  "WHERE name = ?")
	var description string
	var services string
	var id int64
	u := atium.client_serviceInfo{}
	row := ms.db.QueryRow(qs, email)
	
	err := row.Scan(&id,&u.name,&description,&services,&u.created_at,&u.modified_at,
			&u.enabledAt,&u.enabled)
	if err != nil {
		return nil, fmt.Errorf("client_service not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &u.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(u.description)
	}
	err = json.Unmarshal([]byte(services), &u.services)
	if err != nil {
		fmt.Printf("unmarshalling services failed: %v", err)
		fmt.Println(u.services)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//5.medical
func (ms *dbStorage) getmedical(name string) (*atium.medicalInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select M.id,U.name,U.gender,U.contact_id,U.Dob,",
			  "U.created_at,M.medical_record from medical M",
			  "LEFT JOIN user U ON U.id= M.user_id",
			  "WHERE name = ?")
	var medical_record string
	var id int64
	x := atium.medicalInfo{}
	row := ms.db.QueryRow(qs, email)
	
	err := row.Scan(&id,&x.name,&x.gender,&x.contact_id,&x.Dob,
			&x.created_at,&medical_record)
	if err != nil {
		return nil, fmt.Errorf("medical_record not found: %v", err)
	}
	err = json.Unmarshal([]byte(medical_record), &x.medical_record)
	if err != nil {
		fmt.Printf("unmarshalling medical_record failed: %v", err)
		fmt.Println(x.medical_record)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	x.Stats = *stats
	return &x, err
}

//6.survey
func (ms *dbStorage) getsurvey(name string) (*atium.surveyInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select S.id,U.name,U.gender,U.dob,U.created_at,",
			  "S.client_form_map_id,S.form_entry_id from survey S",
			  "LEFT JOIN user U ON U.id= S.user_id",
			  "WHERE name = ?")
	var id int64
	u := atium.surveyInfo{}
	row := ms.db.QueryRow(qs, email)
	
	err := row.Scan(&u.id,&u.name,&u.gender,&u.dob,&u.created_at,
			&u.client_form_map_id,&u.form_entry_id)
	if err != nil {
		return nil, fmt.Errorf("survey not found: %v", err)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//7.stats
func (ms *dbStorage) getstats(name string) (*atium.statsInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
	"select U.id,U.name,U.gender,U.dob,U.created_at, S.entry_at,",
	"S.height,S.weight,S.arms,S.chest,S.waist,S.hips,S.thighs,S.calves from Stats S",
	"LEFT JOIN user U ON U.id = S.id",
	"WHERE name = ?")

	var id int64
	u := atium.statsInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&u.id,&u.name,&u.gender,&u.dob,&u.created_at, &u.entry_at,
	&u.height,&u.weight,&u.arms,&u.chest,&u.waist,&u.hips,&u.thighs,&u.calves )
	if err != nil {
		return nil, fmt.Errorf("stats not found: %v", err)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//8.contact
func (ms *dbStorage) getcontact(email string) (*atium.contactInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
	"select U.id,U.email,U.address,U.primary_ph,U.secondary_ph from user U",
	"WHERE email = ?")

	var address string
	var id int64
	u := atium.contactInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&u.email,&u.address,&u.primary_ph,&u.secondary_ph)
	if err != nil {
		return nil, fmt.Errorf("contact not found: %v", err)
	}
	err = json.Unmarshal([]byte(address), &u.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(u.Address)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//9.Activity
func (ms *dbStorage) getactivity(name string) (*atium.activityInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select A.id,A.name,A.description,A.duration,A.cost from activity A",
			  "WHERE name = ?")

	var description string
	var id int64
	u := atium.activityInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&u.name,&description,&u.duration,&u.cost)
	if err != nil {
		return nil, fmt.Errorf("activity not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &u.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(u.description)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//10.user_activity
func (ms *dbStorage) getuser_activity(name string) (*atium.user_activityInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select A.id,U.name,U.gender,U.contact_id,U.dod,U.created_at,",
			  "A.activity_id,A.user_id,A.start_time,A.end_time from user A",
			  "LEFT JOIN user U on U.id = A.user_id",
			  "WHERE name = ?")

	var id int64
	u := atium.user_activityInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&u.name,&u.gender,&u.contact_id,&u.dod,&u.created_at,
			&u.activity_id,&u.user_id,&u.start_time,&u.end_time)
	if err != nil {
		return nil, fmt.Errorf("user_activity not found: %v", err)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//11.client_form
func (ms *dbStorage) getclient_form(name string) (*atium.client_formInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select F.id,C.name,C.contact_id,F.client_id,C.description,U.services,U.created_at,U.modified_at,",
			  "F.form_name,F.form_label from client_form F ",
			  "LEFT JOIN client C ON C.id = F.client_id",
			  "WHERE name = ?")
	var description string
	var services string
	var id int64
	u := atium.client_formInfo{}
	row := ms.db.QueryRow(qs, email)
	err := row.Scan(&id,&u.name,&u.contact_id,&u.client_id,&description,&services,&u.created_at,&u.modified_at,
			&u.form_name,&u.form_label)
	if err != nil {
	return nil, fmt.Errorf("client_form not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &u.description)
	if err != nil {
	fmt.Printf("unmarshalling description failed: %v", err)
	fmt.Println(u.description)
	}
	err = json.Unmarshal([]byte(services), &u.services)
	if err != nil {
	fmt.Printf("unmarshalling services failed: %v", err)
	fmt.Println(u.services)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
	return nil, fmt.Errorf("err getting stats: %v", err)
}
	u.Stats = *stats
	return &u, err
}

//12.client_activity
func (ms *dbStorage) getclient_activity(name string) (*atium.client_activityInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select C.id,C.client_id,C.activity_id,A.name,",
			  "A.description,A.duration,A.cost from client_activity C",
			  "LEFT JOIN activity A ON A.id = C.activity_id",
			  "WHERE name = ?")
	var description string
	var id int64
	x := atium.client_activityInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&x.client_id,&x.activity_id,&x.name,
			&x.description,&x.duration,&x.cost)
	if err != nil {
		return nil, fmt.Errorf("client_activity not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &x.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(x.description)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	x.Stats = *stats
	return &x, err
}
//13.client_user
func (ms *dbStorage) getclient_user(name string) (*atium.client_userInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select U.id,U.client_id,U.User_id,C.name,C.description,C.contact_id,",
			  "C.service,C.created_at,C.modified_at from client_user U",
			  "LEFT JOIN client C ON C.id = U.client_id",
			  "WHERE name = ?")
	var description string
	var services string
	var id int64
	x := atium.client_userInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&x.client_id,&x.User_id,&x.name,&description,&x.contact_id,
			&service,&x.created_at,&x.modified_at)
	if err != nil {
		return nil, fmt.Errorf("client_user not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &x.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(x.description)
	}
	err = json.Unmarshal([]byte(services), &x.services)
	if err != nil {
		fmt.Printf("unmarshalling services failed: %v", err)
		fmt.Println(x.services)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	x.Stats = *stats
	return &x, err
}

//14.consumption
func (ms *dbStorage) getconsumption(food_id string) (*atium.consumptionInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
		 "select C.id,C.food_id,C.calories_in,C.protein,C.carbs,",
		 "C.fat,C.entry_at from consumption C",
		 "WHERE food_id = ?")
	var id int64
	u := atium.consumptionInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&u.id,&u.food_id,&u.calories_in,&u.protein,&u.carbs,
			&u.fat,&u.entry_at)
	if err != nil {
		return nil, fmt.Errorf("consumption not found: %v", err)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//15.user_consumption
func (ms *dbStorage) getuser_consumption(user_id string) (*atium.user_consumptionInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select U.id,U.user_id,U.consumption_id, C.food_id,C.calories_in,",
			  "C.protein,C.carbs,C.fat,C.entry_at from consumption C",
			  "LEFT JOIN user_consumption U on U.consumption_id = C.id ",
			  "WHERE user_id = ?")
	var id int64
	x := atium.user_consumptionInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&x.user_id,&x.consumption_id,&x.food_id,&x.calories_in,
			&x.protein,&x.carbs,&x.fat,&x.entry_at)
	if err != nil {
		return nil, fmt.Errorf("user_consumption not found: %v", err)
	}

	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	x.Stats = *stats
	return &x, err
}

//16.user_stats
func (ms *dbStorage) getuser_stats(user_id string) (*atium.user_statsInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
			  "select U.id,U.user_id,U.stats_id,S.entry_at,S.height,S.weight,S.arms,S.chest,",
			  "S.waist,S.hips,S.thighs,S.calves from user_stats U",
			  "LEft JOIN stats S ON S.id = U.stats_id",
			  "WHERE user_id = ?")

	var id int64
	x := atium.user_statsInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&x.user_id,&x.stats_id,&x.entry_at,&x.height,&x.weight,&x.arms,&x.chest,
	&x.waist,&x.hips,&x.thighs,&x.calves)
	if err != nil {
		return nil, fmt.Errorf("user_stats not found: %v", err)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	x.Stats = *stats
	return &x, err
}

	  func main(){
  lambda.Start(handle)
}
	  
