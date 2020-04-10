package service

import (
	"atium/pkg/atium"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	asserts "github.com/stretchr/testify/assert"
)

type execResult struct {
	lastInsertId, rowsAffected int64
}

func (er *execResult) LastInsertId() (int64, error) {
	return er.lastInsertId, nil
}
func (er *execResult) RowsAffected() (int64, error) {
	return er.rowsAffected, nil
}

func Test_NewStore(t *testing.T) {
	assert := asserts.New(t)
	s, err := newStore("test", "pwd", "", "db")
	if assert.Nil(err) {
		assert.NotNil(s)
	}
}

func Test_DbStorage_GetUser(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() { _ = db.Close() }()
		dob, _ := time.Parse("2006-Jan-02", "1986-Oct-30")
		createdAt := time.Now()
		rows := sqlMock.NewRows([]string{"id", "name", "dob", "gender",
			"created_at", "address", "primary_ph", "secondary_ph"}).
			AddRow(1, "testuser", dob, "M", createdAt,
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\", \"city\":\"chennai\"}",
				"+21-222222", "")
		qs := fmt.Sprintf("%s %s %s %s",
			"SELECT U.id, U.name, U.dob, U.gender, U.created_at,",
			"C.address, C.primary_ph, C.secondary_ph FROM user U",
			"LEFT JOIN contact C ON C.email = U.email",
			"WHERE U.email = ?")
		mock.ExpectQuery(regexp.QuoteMeta(qs)).
			WithArgs("testuser@atium.com").WillReturnRows(rows)

		rows = sqlMock.NewRows([]string{"entry_at", "height", "weight", "arms",
			"chest", "waist", "hips", "thighs", "calves"}).
			AddRow(createdAt, 178, 75, 0, 0, 0, 0, 0, 0)
		qs = fmt.Sprintf("%s %s %s %s",
			"SELECT entry_at, height, weight, arms, chest,",
			"waist, hips, thighs, calves from stats S",
			"WHERE id = (SELECT stats_id from user_stats",
			"WHERE id = (SELECT MAX(id) FROM user_stats WHERE user_id = ?))")
		mock.ExpectQuery(regexp.QuoteMeta(qs)).
			WithArgs(1).WillReturnRows(rows)
		s := &dbStorage{db}
		u, err := s.getUser("testuser@atium.com")
		if assert.Nil(err) {
			assert.Equal("testuser", u.Name)
			assert.Equal("testuser@atium.com", u.Email)
			assert.Equal(dob, u.Dob)
			assert.Equal("M", u.Gender)
			assert.Equal(createdAt, u.CreatedAt)
			assert.Equal("+21-222222", u.PrimaryPh)
			assert.Equal("", u.SecondaryPh)
			assert.Equal("aaaaa", u.Address.Line1)
			assert.Equal("", u.Address.Line2)
			assert.Equal("6002133", u.Address.Pincode)
			assert.Equal("chennai", u.Address.City)
			assert.Equal("", u.Address.State)
			assert.Equal("IN", u.Address.Country)
			assert.Equal(createdAt, u.Stats.EntryAt)
			assert.Equal(float64(178), u.Stats.Height)
			assert.Equal(float64(75), u.Stats.Weight)
		}
	}
}

func Test_DbStorage_GetClient(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()
	if assert.Nil(err) {
		defer func() { _ = db.Close() }()
		createdAt := time.Now()
		modifiedAt := time.Now()
		rows := sqlMock.NewRows([]string{"id", "description", "Email",
			"created_at", "modifiedAt", "address", "primary_ph", "secondary_ph"}).
			AddRow(1, "AAA", "testuser@atium.com", createdAt, modifiedAt,
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\", \"city\":\"chennai\"}",
				"+21-222222", "")
		qs := fmt.Sprintf("%s %s %s ",
			"SELECT CL.id, CL.description, CL.created_at, CL.modified_at,",
			"C.email, C.address, C.primary_ph, C.secondary_ph FROM client CL",
			"LEFT JOIN contact C ON C.id = CL.contact_id WHERE name = ?")
		mock.ExpectQuery(regexp.QuoteMeta(qs)).
			WithArgs("testuser").WillReturnRows(rows)
		s := &dbStorage{db}
		u, err := s.getClient("testuser")
		if assert.Nil(err) {
			//	assert.Equal(1, u.id)
			assert.Equal("AAA", u.Description)
			assert.Equal("testuser@atium.com", u.Email)
			//assert.Equal("BBB", u.Services)
			assert.Equal(createdAt, u.CreatedAt)
			assert.Equal(modifiedAt, u.ModifiedAt)
			assert.Equal("aaaaa", u.Address.Line1)
			assert.Equal("", u.Address.Line2)
			assert.Equal("6002133", u.Address.Pincode)
			assert.Equal("chennai", u.Address.City)
			assert.Equal("", u.Address.State)
			assert.Equal("IN", u.Address.Country)
			assert.Equal("+21-222222", u.PrimaryPh)
			assert.Equal("", u.SecondaryPh)

		}
	}
}

func Test_DbStorage_UpsertUser(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()
		dob, _ := time.Parse("2006-Jan-02", "1986-Oct-30")

		qs := fmt.Sprintf("%s %s %s %s", "INSERT contact SET email = ?,",
			"address = ?, primary_ph = ?",
			"ON DUPLICATE KEY UPDATE",
			"address = ?, primary_ph = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).
			WithArgs("testuser@atium.com",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"+21-222222",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"+21-222222").
			WillReturnResult(driver.ResultNoRows)

		qs = fmt.Sprintf("%s %s %s %s", "INSERT user SET email = ?,",
			"name = ?, dob = ?, gender = ?",
			"ON DUPLICATE KEY UPDATE",
			"name = ?, dob = ?, gender = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).
			WithArgs("testuser@atium.com", "testuser", dob, "M", "testuser", dob, "M").
			WillReturnResult(driver.ResultNoRows)

		uidRow := sqlMock.NewRows([]string{"id"}).AddRow(22)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM user WHERE email = ?")).
			WithArgs("testuser@atium.com").WillReturnRows(uidRow)

		qs = fmt.Sprintf("%s %s", "INSERT stats SET user_id = ?, height = ?,",
			"weight = ?, arms = ?, chest = ?, waist = ?, hips = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).
			WithArgs(22, 176.0, 75.0, 23.4, 54.5, 34.4, 33.0).WillReturnResult(&execResult{lastInsertId: 1})

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO user_stats (user_id, stats_id) VALUES (?, ?)")).
			WithArgs(22, 1).WillReturnResult(driver.ResultNoRows)
		s := &dbStorage{db}
		err := s.upsertUser(atium.UserDetails{Email: "testuser@atium.com",
			Name: "testuser", Dob: dob, Gender: "M", PrimaryPh: "+21-222222",
			Address: atium.AddressInfo{Line1: "aaaaa", Pincode: "6002133", Country: "IN"},
			Stats: atium.StatsInfo{
				Height: float64(176), Weight: float64(75), Arms: float64(23.4),
				Chest: float64(54.5), Waist: float64(34.4), Hips: float64(33)}})
		assert.Nil(err)
	}
}

func Test_DbStorage_UpsertClient(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()
		qs := fmt.Sprintf("%s %s %s %s", "INSERT contact SET email = ?,",
			"address = ?, primary_ph = ?, secondary_ph = ?",
			"ON DUPLICATE KEY UPDATE",
			"address = ?, primary_ph = ?, secondary_ph = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).
			WithArgs("testclient@atium.com",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"1234", "55555",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"1234", "55555").
			WillReturnResult(driver.ResultNoRows)

		qs = fmt.Sprintf("%s %s %s", "INSERT user SET name = ?, description = ?",
			"ON DUPLICATE KEY UPDATE",
			"modified_at = ? description = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).WillReturnResult(driver.ResultNoRows)

		s := &dbStorage{db}
		err := s.upsertClient(atium.ClientDetails{Email: "testclient@atium.com",
			Name: "testclient", Description: "aaaaa", PrimaryPh: "1234", SecondaryPh: "55555",
			Address: atium.AddressInfo{Line1: "aaaaa", Pincode: "6002133", Country: "IN"}})
		assert.Nil(err)
	}
}

func Test_DbStorage_UpsertClientError1(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()
		qs := fmt.Sprintf("%s %s %s %s", "INSERT contact SET email = ?,",
			"address = ?, primary_ph = ?, secondary_ph = ?",
			"ON DUPLICATE KEY UPDATE",
			"address = ?, primary_ph = ?, secondary_ph = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).
			WithArgs("testclient@atium.com",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"1234", "55555",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"1234", "55555").WillReturnError(fmt.Errorf("no rows found"))

		s := &dbStorage{db}
		err := s.upsertClient(atium.ClientDetails{Email: "testclient@atium.com",
			Name: "testclient", Description: "aaaaa", PrimaryPh: "1234", SecondaryPh: "55555",
			Address: atium.AddressInfo{Line1: "aaaaa", Pincode: "6002133", Country: "IN"}})
		assert.Equal(err.Error(), "no rows found")
	}
}

func Test_DbStorage_UpsertClientError2(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()
		qs := fmt.Sprintf("%s %s %s %s", "INSERT contact SET email = ?,",
			"address = ?, primary_ph = ?, secondary_ph = ?",
			"ON DUPLICATE KEY UPDATE",
			"address = ?, primary_ph = ?, secondary_ph = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).
			WithArgs("testclient@atium.com",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"1234", "55555",
				"{\"line1\":\"aaaaa\",\"pincode\":\"6002133\",\"country\":\"IN\"}",
				"1234", "55555").
			WillReturnResult(driver.ResultNoRows)

		qs = fmt.Sprintf("%s %s %s", "INSERT user SET name = ?, description = ?",
			"ON DUPLICATE KEY UPDATE",
			"modified_at = ? description = ?")
		mock.ExpectExec(regexp.QuoteMeta(qs)).WillReturnError(fmt.Errorf("no rows found"))

		s := &dbStorage{db}
		err := s.upsertClient(atium.ClientDetails{Email: "testclient@atium.com",
			Name: "testclient", Description: "aaaaa", PrimaryPh: "1234", SecondaryPh: "55555",
			Address: atium.AddressInfo{Line1: "aaaaa", Pincode: "6002133", Country: "IN"}})
		assert.Equal(err.Error(), "no rows found")
	}
}

func Test_DbStorage_GetClientUsers(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()
	if assert.Nil(err) {
		defer func() { _ = db.Close() }()
		rows := sqlMock.NewRows([]string{"Email"}).
			AddRow("testuser@atium.com")
		qs := fmt.Sprintf("%s %s",
			"SELECT U.Email FROM user U LEFT JOIN client_user CU ON CU.user_id = U.id",
			"WHERE client_id = ?")
		mock.ExpectQuery(regexp.QuoteMeta(qs)).
			WithArgs("1").WillReturnRows(rows)
		s := &dbStorage{db}
		u, err := s.getClientUsers("1")
		if assert.Nil(err) {
			assert.Equal("testuser@atium.com", u.Email)
		}
	}
}
