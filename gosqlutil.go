package gosqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/tkanos/gonfig"
	"golang.org/x/term"
)

type Configuration struct {
	Host     string
	Port     int
	Database string
	User     string
}

func Conn() *sql.DB {
	// load configuration file
	configuration := Configuration{}
	errc := gonfig.GetConf("./config.json", &configuration)
	if errc != nil {
		log.Fatal("Error loading configuration file: ", errc.Error())
	}

	var connString string
	// check if local subnet or not
	if configuration.Host[0:3] == "192" {
		// get login creds
		fmt.Printf("Enter SQL Server password, user: %s: ", configuration.User)
		bytepw, err1 := term.ReadPassword(int(syscall.Stdin))
		if err1 != nil {
			os.Exit(1)
		}
		pass := string(bytepw)
		fmt.Println("")

		// Build connection string, user configuration file settings
		// Connection with UID + PW
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
			configuration.Host, configuration.User, pass, configuration.Port, configuration.Database)
	} else {
		//Connection with Trusted Credentials
		connString = fmt.Sprintf("server=%s;trusted_connection=yes;database=%s;", configuration.Host, configuration.Database)
	}
	var err error

	// Create connection pool
	var db *sql.DB
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db

}
