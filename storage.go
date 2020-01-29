package service

import (
	"atium/pkg/atium"
	"database/sql"
	"encoding/json"
	"fmt"
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

func handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Event: ", event)
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
func (ms *dbStorage) Close() error {
	return ms.db.Close()
}

// -----------------------------------------------------------------------------
// Client related DB function which implements interface 'storage'
// -----------------------------------------------------------------------------

func (ms *dbStorage) listClients() (interface{}, error) {
	return nil, nil
}

func (ms *dbStorage) getClient(c string) (*atium.ClientInfo, error) {
	return &atium.ClientInfo{}, nil
}

func (ms *dbStorage) upsertClient(q interface{}) error {
	return nil
}

func (ms *dbStorage) deleteClient(a string) error {
	return nil
}

// -----------------------------------------------------------------------------
// User related DB function which implements interface 'storage'
// -----------------------------------------------------------------------------

func (ms *dbStorage) listUsers() (interface{}, error) {
	return nil, nil
}

func (ms *dbStorage) getUser(email string) (*atium.UserInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
		"SELECT U.id, U.name, C.email, U.dob, U.created_at, C.address,",
		"C.primary_ph, C.secondary_ph from user U",
		"LEFT JOIN contact C ON C.id = U.contact_id",
		"WHERE email = ?")

	var address string
	var id int64
	u := atium.UserInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id, &u.Name, &u.Email, &u.Dob, &u.CreatedAt,
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
	u.Stats = stats
	return &u, err
}

func (ms *dbStorage) upsertUser(q interface{}) error {
	return nil
}

func (ms *dbStorage) deleteUser(a string) error {
	return nil
}

func main() {
	lambda.Start(handle)
}
