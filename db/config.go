package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var sqlDb *sql.DB

type Dbconfig struct {
	hostname string
	password string
	dbname   string
	username string
	port     string
	dbserver string
}

func (dbc *Dbconfig) Dbname() string {
	return dbc.dbname
}

func (dbc *Dbconfig) SetDbname(dbname string) {
	dbc.dbname = dbname
}

func (dbc *Dbconfig) Hostname() string {
	return dbc.hostname
}

func (dbc *Dbconfig) SetHostname(hname string) {
	dbc.hostname = hname
}

func (dbc *Dbconfig) Username() string {
	return dbc.username
}

func (dbc *Dbconfig) SetUsername(userid string) {
	dbc.username = userid
}

func (dbc *Dbconfig) Password() string {
	return dbc.password
}

func (dbc *Dbconfig) SetPassword(password string) {
	dbc.password = password
}

func (dbc *Dbconfig) Dbserver() string {
	return dbc.dbserver
}

func (dbc *Dbconfig) SetDbserver(server string) {
	dbc.dbserver = server
}

func (dbc *Dbconfig) Port() string {
	return dbc.port
}

func (dbc *Dbconfig) SetPort(port string) {
	dbc.port = port
}

func Connect(dbc *Dbconfig) *sql.DB {

	server := dbc.Dbserver()

	cfg := mysql.Config{
		User:                 dbc.Username(),
		Passwd:               dbc.Password(),
		Net:                  "tcp",
		Addr:                 dbc.Hostname() + ":" + dbc.Port(),
		DBName:               dbc.Dbname(),
		AllowNativePasswords: true,
	}
	switch server {
	case "mysql":
		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			panic(err)
		}

		//defer db.Close()

		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}
		// See "Important settings" section.
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		fmt.Println("connected")

		return db
	}

	return nil

}
