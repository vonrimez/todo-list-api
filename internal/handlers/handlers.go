package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vonrimez/TaskAPI/domain"
	"github.com/vonrimez/TaskAPI/internal/models"
	"github.com/vonrimez/TaskAPI/internal/service"
)

type H map[string]any
type contextKey string

var (
	UserID contextKey = "user_id"
	Error  contextKey = "error"
)

type Handler struct {
	taskService *service.TaskService
	userService *service.UserService
	validate    *service.Validate
}

func NewHandler(task *service.TaskService, user *service.UserService, val *service.Validate) *Handler {
	return &Handler{
		taskService: task, userService: user, validate: val,
	}
}

func getUserID(r *http.Request) (int, error) {
	rawID := getValue(r, UserID)
	id, isInt := rawID.(int)
	if !isInt {
		return 0, domain.NewBadRequestError("invalid parameter \"ID\"")
	}
	return id, nil
}

func getTaskID(r *http.Request) (int, error) {
	var rawID string
	chi.URLParam(r, "id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return 0, domain.NewBadRequestError("ID must be an integer")
	}
	return id, nil
}

func (hdl *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	inputUser := models.UserLoginInput{}
	err := hdl.validate.ShouldBindJSON(r.Body, &inputUser)
	if err != nil {
		fieldErrors := hdl.validate.ValError(err)
		if fieldErrors != nil {
			respondWithJSON(w, &Response{
				statusCode: http.StatusBadRequest,
				body:       H{"error": fieldErrors},
			})
			setValue(r, Error, err)
			return
		}
		respondWithJSON(w, &Response{
			statusCode: http.StatusBadRequest,
			body:       H{"error": "form is incorrect"},
		})
		setValue(r, Error, err)
		return
	}

	outputUser, err := hdl.userService.Login(inputUser)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}

	respondWithJSON(w, &Response{
		statusCode: http.StatusOK,
		body: H{
			"access_token": outputUser.JWT,
			"token_type":   "Bearer",
			"user":         *outputUser,
		},
	})
}

func (hdl *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	inputUser := models.UserRegisterInput{}
	err := hdl.validate.ShouldBindJSON(r.Body, &inputUser)
	if err != nil {
		fieldErrors := hdl.validate.ValError(err)
		if fieldErrors != nil {
			respondWithJSON(w, &Response{
				statusCode: http.StatusBadRequest,
				body:       H{"error": fieldErrors},
			})
			setValue(r, Error, err)
			return
		}
		respondWithJSON(w, &Response{
			statusCode: http.StatusBadRequest,
			body:       H{"error": "form is incorrect"},
		})
		setValue(r, Error, err)
		return
	}

	outputUser, err := hdl.userService.Register(inputUser)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}

	respondWithJSON(w, &Response{
		statusCode: http.StatusCreated,
		body: H{
			"access_token": outputUser.JWT,
			"token_type":   "Bearer",
			"user":         *outputUser,
		},
	})
}

func (hdl *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	outputTasks, err := hdl.taskService.GetAll(userID)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	respondWithJSON(w, &Response{
		statusCode: http.StatusOK,
		body:       outputTasks,
	})
}

func (hdl *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	userID, err := getUserID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	outputTask, err := hdl.taskService.GetById(taskID, userID)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	respondWithJSON(w, &Response{
		statusCode: http.StatusOK,
		body:       outputTask,
	})
}

func (hdl *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	userID, err := getUserID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	outputTask, err := hdl.taskService.Delete(taskID, userID)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	respondWithJSON(w, &Response{
		statusCode: http.StatusOK,
		body:       H{"deleted_task": *outputTask},
	})
}

func (hdl *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}

	inputTask := models.TaskUpdateInput{}
	err = hdl.validate.ShouldBindJSON(r.Body, &inputTask)
	if err != nil {
		fieldErrors := hdl.validate.ValError(err)
		if fieldErrors != nil {
			respondWithJSON(w, &Response{
				statusCode: http.StatusBadRequest,
				body:       H{"error": fieldErrors},
			})
			setValue(r, Error, err)
			return
		}
		respondWithJSON(w, &Response{
			statusCode: http.StatusBadRequest,
			body:       H{"error": "form is incorrect"},
		})
		setValue(r, Error, err)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	outputTask, err := hdl.taskService.Update(inputTask, taskID, userID)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	respondWithJSON(w, &Response{
		statusCode: http.StatusOK,
		body:       H{"updated_task": *outputTask},
	})
}

func (hdl *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	inputTask := models.TaskCreateInput{}
	err := hdl.validate.ShouldBindJSON(r.Body, &inputTask)
	if err != nil {
		fieldErrors := hdl.validate.ValError(err)
		if fieldErrors != nil {
			respondWithJSON(w, &Response{
				statusCode: http.StatusBadRequest,
				body:       H{"error": fieldErrors},
			})
			setValue(r, Error, err)
			return
		}
		respondWithJSON(w, &Response{
			statusCode: http.StatusBadRequest,
			body:       H{"error": "form is incorrect"},
		})
		setValue(r, Error, err)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}

	outputTask, err := hdl.taskService.Create(inputTask, userID)
	if err != nil {
		respondWithAppError(w, err)
		setValue(r, Error, err)
		return
	}
	respondWithJSON(w, &Response{
		statusCode: http.StatusCreated,
		body:       H{"created_task": *outputTask},
	})
}
