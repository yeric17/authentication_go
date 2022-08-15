package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/yeric17/thullo/pkg/config"
)

func GetConnection() *sql.DB {

	db, err := sql.Open(config.DB_DRIVER, config.CONNECTION_STRING)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected Database")
	return db
}
