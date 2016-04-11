package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init(host, usr, pwd, dbname string) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", usr, pwd, host, dbname))
	if err != nil {
		panic(err)
	}
}

func Get() *sql.DB {
	return db
}
