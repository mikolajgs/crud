package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "os"
	_ "time"

	restapi "github.com/mikolajgs/prototyping/pkg/rest-api"
	stdb "github.com/mikolajgs/prototyping/pkg/struct-db-postgres"
	"github.com/mikolajgs/prototyping/pkg/ui"

	_ "github.com/lib/pq"
)

const dbUser = "uiuser"
const dbPass = "uipass"
const dbName = "uidb"
const dbPort = "54320"

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost user=%s password=%s port=%s dbname=%s sslmode=disable", dbUser, dbPass, dbPort, dbName))
	if err != nil {
		log.Fatal("Error connecting to db")
	}

	ctl := ui.NewController(db, "ui_")
	apiCtl := restapi.NewController(db, "ui_")
	s2db := stdb.NewController(db, "ui_", nil)

	person := &Person{}
	group := &Group{}

	err2 := s2db.DropTables(person, group)
	if err2 != nil {
		log.Printf("Error with dropping tables: %s", err2.Error())
	}

	err2 = s2db.CreateTables(person, group)
	if err2 != nil {
		log.Fatalf("Error with creating tables: %s", err.Error())
	}
	
	http.Handle("/ui/v1/", ctl.GetHTTPHandler(
		"/ui/v1/",
		func() interface{}{ return &Person{} },
		func() interface{}{ return &Group{} },
	))
	http.Handle("/api/v1/persons/", apiCtl.GetHTTPHandler(
		"/api/v1/persons/",
		func() interface{}{ return &Person{} },
		func() interface{}{ return &Person{} },
		func() interface{}{ return &Person{} },
		func() interface{}{ return &Person{} },
		func() interface{}{ return &Person{} },
		func() interface{}{ return &Person{} },
	))
	http.Handle("/api/v1/groups/", apiCtl.GetHTTPHandler(
		"/api/v1/groups/",
		func() interface{}{ return &Group{} },
		func() interface{}{ return &Group{} },
		func() interface{}{ return &Group{} },
		func() interface{}{ return &Group{} },
		func() interface{}{ return &Group{} },
		func() interface{}{ return &Group{} },
	))
	log.Fatal(http.ListenAndServe(":9001", nil))
}
