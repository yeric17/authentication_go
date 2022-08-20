package data

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/yeric17/thullo/pkg/config"
)

var (
	Connection *sql.DB
)

func init() {
	Connection = GetConnection()
}

func GetConnection() *sql.DB {

	db, err := sql.Open(config.DB_DRIVER, config.CONNECTION_STRING)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)
	fmt.Println("Successfully connected Database")
	return db
}
