package struct2db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

// Global vars used across all the tests
var dbUser = "gocrudtest"
var dbPass = "secret"
var dbName = "gocrud"
var dbConn *sql.DB

var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

var testController *Controller
var testStructObj *TestStruct

// Test struct for all the tests
type TestStruct struct {
	ID    int64 `json:"test_struct_id"`
	Flags int64 `json:"test_struct_flags"`

	// Test email validation
	PrimaryEmail   string `json:"email" validation:"req"`
	EmailSecondary string `json:"email2" validation:"req email"`

	// Test length validation
	FirstName string `json:"first_name" validation:"req lenmin:2 lenmax:30"`
	LastName  string `json:"last_name" validation:"req lenmin:0 lenmax:255"`

	// Test int value validation
	Age   int `json:"age" validation:"valmin:1 valmax:120"`
	Price int `json:"price" validation:"valmin:0 valmax:999"`

	// Test regular expression
	PostCode  string `json:"post_code" validation:"req lenmin:6 regexp:^[0-9]{2}\\-[0-9]{3}$"`
	PostCode2 string `json:"post_code2" validation:"lenmin:6" validation_regexp:"^[0-9]{2}\\-[0-9]{3}$"`

	// Some other fields
	Password        string `json:"password"`
	CreatedByUserID int64  `json:"created_by_user_id"`

	// Test unique tag
	Key string `json:"key" validation:"req uniq lenmin:30 lenmax:255"`
}

func TestMain(m *testing.M) {
	createDocker()
	createController()

	code := m.Run()
	removeDocker()
	os.Exit(code)
}

func createDocker() {
	var err error
	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dockerResource, err = dockerPool.Run("postgres", "13", []string{"POSTGRES_PASSWORD=" + dbPass, "POSTGRES_USER=" + dbUser, "POSTGRES_DB=" + dbName})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	if err = dockerPool.Retry(func() error {
		var err error
		dbConn, err = sql.Open("postgres", fmt.Sprintf("host=localhost user=%s password=%s port=%s dbname=%s sslmode=disable", dbUser, dbPass, dockerResource.GetPort("5432/tcp"), dbName))
		if err != nil {
			return err
		}
		return dbConn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func createController() {
	testController = NewController(dbConn, "gen64_")
	testStructObj = &TestStruct{}
}

func removeDocker() {
	dockerPool.Purge(dockerResource)
}

func getTableNameCnt(tblName string) (int64, error) {
	var cnt int64
	err := dbConn.QueryRow("SELECT COUNT(table_name) AS c FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", tblName).Scan(&cnt)
	return cnt, err
}

func getRow() (int64, int64, string, string, string, string, int, int, string, string, string, int64, string, error) {
	var id, flags, createdByUserID int64
	var primaryEmail, emailSecondary, firstName, lastName, postCode, postCode2, password, key string
	var age, price int
	err := dbConn.QueryRow("SELECT * FROM gen64_test_structs ORDER BY test_struct_id DESC LIMIT 1").Scan(&id, &flags, &primaryEmail, &emailSecondary, &firstName, &lastName, &age, &price, &postCode, &postCode2, &password, &createdByUserID, &key)
	return id, flags, primaryEmail, emailSecondary, firstName, lastName, age, price, postCode, postCode2, password, createdByUserID, key, err
}

func getRowById(id int64) (int64, string, string, string, string, int, int, string, string, string, int64, string, error) {
	var id2, flags, createdByUserID int64
	var primaryEmail, emailSecondary, firstName, lastName, postCode, postCode2, password, key string
	var age, price int
	err := dbConn.QueryRow(fmt.Sprintf("SELECT * FROM gen64_test_structs WHERE test_struct_id = %d", id)).Scan(&id2, &flags, &primaryEmail, &emailSecondary, &firstName, &lastName, &age, &price, &postCode, &postCode2, &password, &createdByUserID, &key)
	return flags, primaryEmail, emailSecondary, firstName, lastName, age, price, postCode, postCode2, password, createdByUserID, key, err
}

func getRowCntById(id int64) (int64, error) {
	var cnt int64
	err := dbConn.QueryRow(fmt.Sprintf("SELECT COUNT(*) AS c FROM gen64_test_structs WHERE test_struct_id = %d", id)).Scan(&cnt)
	return cnt, err
}

func getTestStructWithData() *TestStruct {
	ts := &TestStruct{}
	ts.Flags = 4
	ts.PrimaryEmail = "primary@gen64.net"
	ts.EmailSecondary = "secondary@gen64.net"
	ts.FirstName = "John"
	ts.LastName = "Smith"
	ts.Age = 37
	ts.Price = 444
	ts.PostCode = "00-000"
	ts.PostCode2 = "11-111"
	ts.Password = "yyy"
	ts.CreatedByUserID = 4
	ts.Key = fmt.Sprintf("12345679012345678901234567890%d", time.Now().UnixNano())
	return ts
}

func areTestStructObjectSame(ts1 *TestStruct, ts2 *TestStruct) bool {
	if ts1.Flags != ts2.Flags {
		return false
	}
	if ts1.PrimaryEmail != ts2.PrimaryEmail {
		return false
	}
	if ts1.EmailSecondary != ts2.EmailSecondary {
		return false
	}
	if ts1.FirstName != ts2.FirstName {
		return false
	}
	if ts1.LastName != ts2.LastName {
		return false
	}
	if ts1.Age != ts2.Age {
		return false
	}
	if ts1.Price != ts2.Price {
		return false
	}
	if ts1.PostCode != ts2.PostCode {
		return false
	}
	if ts1.PostCode2 != ts2.PostCode2 {
		return false
	}
	if ts1.Password != ts2.Password {
		return false
	}
	if ts1.CreatedByUserID != ts2.CreatedByUserID {
		return false
	}
	if ts1.Key != ts2.Key {
		return false
	}
	return true
}

