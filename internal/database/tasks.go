package database

import (
	"database/sql"
	"errors"

	"github.com/vonrimez/TaskAPI/domain"
	"github.com/vonrimez/TaskAPI/internal/models"
)

type TasksDB struct {
	db *sql.DB
}

func GetNewTasksDB(db *sql.DB) *TasksDB {
	return &TasksDB{db: db}
}

// need return values: id, title, description, status, created_at, updated_at
func (tdb *TasksDB) execWithError(query string, args ...any) (*models.Task, *domain.AppError) {
	row := tdb.db.QueryRow(query, args...)
	outputTask := models.Task{}
	err := row.Scan(
		&outputTask.ID, &outputTask.Title, &outputTask.Description, &outputTask.Status, &outputTask.Created_at, &outputTask.Updated_at,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("there are no matches")
		}
		return nil, domain.NewInternalError(err)
	}
	return &outputTask, nil
}

func (tdb *TasksDB) GetTasks(userID int) ([]models.Task, *domain.AppError) {
	query := `
	SELECT id, title, description, status, TO_CHAR(created_at, 'DD.MM.YYY HH24:MI:SS'), TO_CHAR(updated_at, 'DD.MM.YYY HH24:MI:SS') 
	FROM tasks
	WHERE user_id = $1;
	`
	rows, err := tdb.db.Query(query, userID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0, 8)

	for rows.Next() {
		t := models.Task{}

		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at,
		)
		if err != nil {
			return nil, domain.NewInternalError(err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return tasks, nil
}

func (tdb *TasksDB) GetTaskById(taskID int, userID int) (*models.Task, *domain.AppError) {
	query := `
	SELECT id, title, description, status, TO_CHAR(created_at, 'DD.MM.YYY HH24:MI:SS'), TO_CHAR(updated_at, 'DD.MM.YYY HH24:MI:SS') 
	FROM tasks 
	WHERE id = $1
	AND user_id = $2;
	`
	row := tdb.db.QueryRow(query, taskID, userID)
	if err := row.Err(); err != nil {
		return nil, domain.NewInternalError(err)
	}

	outputTask := models.Task{}

	err := row.Scan(
		&outputTask.ID, &outputTask.Title, &outputTask.Description, &outputTask.Status, &outputTask.Created_at, &outputTask.Updated_at,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("there are no matches")
		}
		return nil, domain.NewInternalError(err)
	}

	return &outputTask, nil
}

func (tdb *TasksDB) CreateTask(inputTask models.TaskCreateInput, userID int) (*models.Task, *domain.AppError) {
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

	return tdb.execWithError(query, inputTask.Title, inputTask.Description, inputTask.Status, userID)
}

func (tdb *TasksDB) UpdateTask(inputTask models.TaskUpdateInput, taskID int, userID int) (*models.Task, *domain.AppError) {
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

	return tdb.execWithError(query, inputTask.Title, inputTask.Description, inputTask.Status, taskID, userID)
}

func (tdb *TasksDB) DeleteTask(taskID int, userID int) (*models.Task, *domain.AppError) {
	query := `
	DELETE FROM tasks 
	WHERE id = $1 AND user_id = $2 
	RETURNING id, title, description, status, created_at, updated_at;
	`
	return tdb.execWithError(query, taskID, userID)
}
