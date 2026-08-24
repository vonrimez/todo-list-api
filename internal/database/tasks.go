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

// need return values: id, title, description, status, created_at, updated_at
func (tdb *TasksDB) execWithError(query string, args ...any) (*models.Task, error) {
	row := tdb.db.QueryRow(query, args...)
	t := models.Task{}
	err := row.Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (tdb *TasksDB) GetTaskInfo(t *models.Task) string {
	return fmt.Sprintf("[ID:%d]=%s (%s) => %s [%v : %v]\n", t.ID, t.Title, t.Description, t.Status, t.Created_at, t.Updated_at)
}

func (tdb *TasksDB) GetTasks(userID int) ([]models.TaskGetOutput, error) {
	query := `
	SELECT id, title, description, status, TO_CHAR(created_at, 'DD.MM.YYY HH24:MI:SS'), TO_CHAR(updated_at, 'DD.MM.YYY HH24:MI:SS') 
	FROM tasks;
	`
	rows, err := tdb.db.Query(query)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.TaskGetOutput{}

	for rows.Next() {
		t := models.TaskGetOutput{}

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

func (tdb *TasksDB) GetTaskById(taskID int, userID int) (models.TaskGetOutput, error) {
	query := `
	SELECT id, title, description, status, TO_CHAR(created_at, 'DD.MM.YYY HH24:MI:SS'), TO_CHAR(updated_at, 'DD.MM.YYY HH24:MI:SS') 
	FROM tasks 
	WHERE id = $1;
	`
	row := tdb.db.QueryRow(query, taskID)
	if err := row.Err(); err != nil {
		return models.TaskGetOutput{}, err
	}

	t := models.TaskGetOutput{}

	err := row.Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at,
	)
	if err != nil {
		return models.TaskGetOutput{}, nil
	}

	return t, nil
}

func (tdb *TasksDB) CreateTask(task models.TaskCreateInput, userID int) (*models.Task, error) {
	query := `
	INSERT INTO tasks (title, description, status, user_id) 
	VALUES (
		$1, 
		$2, 
		COALESCE(NULLIF($3, '')::TASK_STATUS, 'todo'::TASK_STATUS),
		$4
	) 
	RETURNING id, title, description, status, created_at, updated_at;
	`

	return tdb.execWithError(query, task.Title, task.Description, task.Status, userID)
}

func (tdb *TasksDB) UpdateTask(task models.TaskUpdateInput, taskID int, userID int) (*models.Task, error) {
	query := `
	UPDATE tasks 
	SET 
		title = COALESCE(NULLIF($1, ''), title), 
		description = COALESCE(NULLIF($2, ''), description), 
		status = COALESCE(NULLIF($3, '')::TASK_STATUS, status),
		updated_at = NOW() 
	WHERE id = $4 AND user_id = $5 
	RETURNING id, title, description, status, created_at, updated_at;
	`

	return tdb.execWithError(query, task.Title, task.Description, task.Status, taskID, userID)
}

func (tdb *TasksDB) DeleteTask(taskID int, userID int) (*models.Task, error) {
	query := `
	DELETE FROM tasks 
	WHERE id = $1 AND user_id = $2 
	RETURNING id, title, description, status, created_at, updated_at;
	`
	return tdb.execWithError(query, taskID, userID)
}
