package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("error: failed connection"))
	}

	res, err := db.Exec("INSERT INTO tasks (title, description) VALUES ($1, $2)", "some_title", "some_desc")
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		panic("error: no one rows insert")
	}
	if err != nil {
		panic(err)
	}

	var tasks struct {
		id          int
		title       string
		description string
		status      string
		created_at  time.Time
		updated_at  time.Time
	}

	err = db.QueryRow("SELECT * FROM tasks WHERE id = $1", 1).Scan(
		&tasks.id, &tasks.title, &tasks.description, &tasks.status, &tasks.created_at, &tasks.updated_at,
	)
	if err != nil {
		panic(err)
	}
}
