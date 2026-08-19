package database

import (
	"database/sql"
	"fmt"
)

func EstablishConnection(dbDriver string, dbUrl string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error: could not connect with databse")
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)

	return db, nil
}
