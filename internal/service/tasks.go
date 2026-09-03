package service

import (
	"github.com/vonrimez/TaskAPI/domain"
	"github.com/vonrimez/TaskAPI/internal/models"
)

type TaskRepository interface {
	GetTasks(int) ([]models.Task, error)
	GetTaskById(int, int) (*models.Task, error)
	CreateTask(models.TaskCreateInput, int) (*models.Task, error)
	UpdateTask(models.TaskUpdateInput, int, int) (*models.Task, error)
	DeleteTask(int, int) (*models.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func isNegative(a int) bool {
	if a >= 0 {
		return false
	}
	return true
}

func isNotValidStatus(status string) bool {
	switch status {
	case "todo":
		fallthrough
	case "in-progress":
		fallthrough
	case "done":
		return false
	}
	return true
}

func (s *TaskService) GetAll(userID int) ([]models.Task, error) {
	if isNegative(userID) {
		return nil, domain.NewBadRequestError("the user id must be positive")
	}
	outputTasks, appErr := s.repo.GetTasks(userID)
	if appErr != nil {
		return nil, appErr
	}
	return outputTasks, nil
}

func (s *TaskService) GetById(taskID int, userID int) (*models.Task, error) {
	if isNegative(userID) {
		return nil, domain.NewBadRequestError("the user id must be positive")
	}
	if isNegative(taskID) {
		return nil, domain.NewBadRequestError("the task id must be positive")
	}
	outputTask, appErr := s.repo.GetTaskById(taskID, userID)
	if appErr != nil {
		return nil, appErr
	}
	return outputTask, nil
}

func (s *TaskService) Create(inputTask models.TaskCreateInput, userID int) (*models.Task, error) {
	if isNegative(userID) {
		return nil, domain.NewBadRequestError("the user id must be positive")
	}
	if len(*inputTask.Title) == 0 {
		return nil, domain.NewBadRequestError("the title cannot be empty")
	}
	if inputTask.Status != nil && isNotValidStatus(*inputTask.Status) {
		return nil, domain.NewBadRequestError("the status must be 'todo', 'in-progress' or 'done'")
	}
	outputTask, appErr := s.repo.CreateTask(inputTask, userID)
	if appErr != nil {
		return nil, appErr
	}
	return outputTask, nil
}

func (s *TaskService) Update(inputTask models.TaskUpdateInput, taskID int, userID int) (*models.Task, error) {
	if isNegative(userID) {
		return nil, domain.NewBadRequestError("the user id must be positive")
	}
	if isNegative(taskID) {
		return nil, domain.NewBadRequestError("the task id must be positive")
	}
	if inputTask.Title == nil &&
		inputTask.Description == nil &&
		inputTask.Status == nil {
		return nil, domain.NewBadRequestError("there's nothing to update.")
	}
	if inputTask.Status != nil && isNotValidStatus(*inputTask.Status) {
		return nil, domain.NewBadRequestError("the status must be 'todo', 'in-progress' or 'done'")
	}
	outputTask, appErr := s.repo.UpdateTask(inputTask, taskID, userID)
	if appErr != nil {
		return nil, appErr
	}
	return outputTask, nil
}

func (s *TaskService) Delete(taskID int, userID int) (*models.Task, error) {
	if isNegative(userID) {
		return nil, domain.NewBadRequestError("the user id must be positive")
	}
	if isNegative(taskID) {
		return nil, domain.NewBadRequestError("the task id must be positive")
	}
	outputTask, appErr := s.repo.DeleteTask(taskID, userID)
	if appErr != nil {
		return nil, appErr
	}
	return outputTask, nil
}
