package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/domain"
	"github.com/vonrimez/TaskAPI/internal/models"
	"github.com/vonrimez/TaskAPI/internal/service"
)

type Handler struct {
	taskService *service.TaskService
	userService *service.UserService
}

func GetNewHandler(task *service.TaskService, user *service.UserService) *Handler {
	return &Handler{
		taskService: task, userService: user,
	}
}

func getUserID(context *gin.Context) (int, error) {
	rawID, exists := context.Get("ID")
	if !exists {
		return 0, domain.NewBadRequestError("cannot get param for key \"ID\"")
	}
	id, isInt := rawID.(int)
	if !isInt {
		return 0, domain.NewBadRequestError("param for key \"ID\" should be an integer")
	}
	return id, nil
}

func getTaskID(context *gin.Context) (int, error) {
	rawID := context.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return 0, domain.NewBadRequestError("ID must be an integer")
	}
	return id, nil
}

func respondWithError(context *gin.Context, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		context.JSON(appErr.Code, appErr.Message)
		context.Error(appErr.Err)
		return
	}
	context.JSON(http.StatusInternalServerError, "internal error")
	context.Error(err)
}

func (hdl *Handler) UserLogin(context *gin.Context) {
	inputUser := models.UserLoginInput{}
	err := context.ShouldBindJSON(&inputUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, "form is incorrect")
		context.Error(err)
		return
	}

	outputUser, err := hdl.userService.Login(inputUser)
	if err != nil {
		respondWithError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token": outputUser.JWT,
		"token_type":   "Bearer",
		"user":         &outputUser,
	})
}

func (hdl *Handler) UserRegister(context *gin.Context) {
	inputUser := models.UserRegisterInput{}
	err := context.ShouldBindJSON(&inputUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, "form is incorrect")
		context.Error(err)
		return
	}

	outputUser, err := hdl.userService.Register(inputUser)
	if err != nil {
		respondWithError(context, err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"access_token": outputUser.JWT,
		"token_type":   "Bearer",
		"user":         &outputUser,
	})
}

func (hdl *Handler) GetTasks(context *gin.Context) {
	userID, err := getUserID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}
	outputTasks, err := hdl.taskService.GetAll(userID)
	if err != nil {
		respondWithError(context, err)
		return
	}
	context.JSON(http.StatusOK, outputTasks)
}

func (hdl *Handler) GetTaskById(context *gin.Context) {
	taskID, err := getTaskID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}
	userID, err := getUserID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}
	outputTask, err := hdl.taskService.GetById(taskID, userID)
	if err != nil {
		respondWithError(context, err)
		return
	}
	context.JSON(http.StatusOK, outputTask)
}

func (hdl *Handler) DeleteTask(context *gin.Context) {
	taskID, err := getTaskID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}
	userID, err := getUserID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}
	outputTask, err := hdl.taskService.Delete(taskID, userID)
	if err != nil {
		respondWithError(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"deleted_task": outputTask})
}

func (hdl *Handler) UpdateTask(context *gin.Context) {
	taskID, err := getTaskID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}

	inputTask := models.TaskUpdateInput{}
	err = context.ShouldBindJSON(&inputTask)
	if err != nil {
		context.JSON(http.StatusBadRequest, "form is incorrect")
		context.Error(err)
		return
	}
	userID, err := getUserID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}
	outputTask, err := hdl.taskService.Update(inputTask, taskID, userID)
	if err != nil {
		respondWithError(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"updated_task": outputTask})
}

func (hdl *Handler) CreateTask(context *gin.Context) {
	inputTask := models.TaskCreateInput{}
	err := context.ShouldBindJSON(&inputTask)
	if err != nil {
		context.JSON(http.StatusBadRequest, "form is incorrect")
		context.Error(err)
		return
	}

	userID, err := getUserID(context)
	if err != nil {
		respondWithError(context, err)
		return
	}

	outputTask, err := hdl.taskService.Create(inputTask, userID)
	if err != nil {
		respondWithError(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"created_task": outputTask})
}
