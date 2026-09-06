package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func EstablishConnection(dbDriver string, dbUrl string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err // fmt.Errorf("error: could not connect with databse")
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)

	return db, nil
}
