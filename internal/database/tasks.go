package database

import "database/sql"

type TasksDB struct {
	db *sql.DB
}

func GetNewTasksDB(db *sql.DB) *TasksDB {
	return &TasksDB{db: db}
}
