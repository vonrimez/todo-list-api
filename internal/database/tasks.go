package database

import (
	"database/sql"
	"fmt"

	"github.com/vonrimez/TaskAPI/internal/models"
)

type TasksDB struct {
	db *sql.DB
}

func GetNewTasksDB(db *sql.DB) *TasksDB {
	return &TasksDB{db: db}
}

func (tdb *TasksDB) execWithError(query string, args ...any) error {
	res, err := tdb.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rA, err := res.RowsAffected()
	if rA == 0 || err != nil {
		return fmt.Errorf("error: no action was taken on the database")
	}
	return nil
}

func (tdb *TasksDB) GetTasks(userID int) ([]models.Task, error) {
	query := `SELECT * FROM tasks;`
	rows, err := tdb.db.Query(query)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		t := models.Task{}

		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (tdb *TasksDB) GetTaskById(taskID int, userID int) (models.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1;`
	row := tdb.db.QueryRow(query, taskID)
	if err := row.Err(); err != nil {
		return models.Task{}, err
	}

	t := models.Task{}

	err := row.Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at,
	)
	if err != nil {
		return models.Task{}, nil
	}

	return t, nil
}

func (tdb *TasksDB) CreateTask(task models.TaskCreateInput, userID int) error {
	query := `
	INSERT INTO tasks (title, description, status) 
	VALUES (
		$1, 
		$2, 
		COALESCE(NULLIF($3, '')::TASK_STATUS, 'todo'::TASK_STATUS)
	);
	`

	return tdb.execWithError(query, task.Title, task.Description, task.Status)
}

func (tdb *TasksDB) UpdateTask(task models.TaskUpdateInput, taskID int, userID int) error {
	query := `
	UPDATE tasks 
	SET 
		title = COALESCE(NULLIF($1, ''), title), 
		description = COALESCE(NULLIF($2, ''), description), 
		status = COALESCE(NULLIF($3, '')::TASK_STATUS, status),
		updated_at = NOW()
	WHERE id = $4;
	`

	return tdb.execWithError(query, task.Title, task.Description, task.Status, taskID)
}

func (tdb *TasksDB) DeleteTask(taskID int, userID int) error {
	query := `DELETE FROM tasks WHERE id = $1;`
	return tdb.execWithError(query, taskID)
}
