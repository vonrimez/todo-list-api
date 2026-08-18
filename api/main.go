package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	res, err := db.Exec("INSERT INTO tasks (title, description) VALUES ($1, $2)", "some_title", "some_description")
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		panic("error: no one rows inserted")
	}
	log.Println("Message was successfully sent")

	var tasks struct {
		ID          int
		Title       string
		Description string
		Status      string
		Created_at  time.Time
		Updated_at  time.Time
	}

	err = db.QueryRow("SELECT * FROM tasks WHERE id = $1", 1).Scan(
		&tasks.ID, &tasks.Title, &tasks.Description, &tasks.Status, &tasks.Created_at, &tasks.Updated_at,
	)
	if err != nil {
		panic(err)
	}
	log.Println(tasks)
}
