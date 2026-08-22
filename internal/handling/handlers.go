package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/internal/database"
	"github.com/vonrimez/TaskAPI/internal/models"
)

type Handler struct {
	tasksdb *database.TasksDB
}

func GetNewHandler(tasksdb *database.TasksDB) *Handler {
	return &Handler{tasksdb: tasksdb}
}

func (hdl *Handler) UserLogin(context *gin.Context) {

}

func (hdl *Handler) UserRegister(context *gin.Context) {

}

func (hdl *Handler) GetTasks(context *gin.Context) {
	tasks, err := hdl.tasksdb.GetTasks()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, tasks)
}

func (hdl *Handler) GetTaskById(context *gin.Context) {
	rawID := context.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		context.String(http.StatusBadRequest, "ID must be an integer")
		return
	}
	task, err := hdl.tasksdb.GetTaskById(id)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, task)
}

func (hdl *Handler) DeleteTask(context *gin.Context) {
	rawID := context.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		context.String(http.StatusBadRequest, "ID must be an integer")
		return
	}

	err = hdl.tasksdb.DeleteTask(id)
	if errors.Is(err, sql.ErrNoRows) {
		context.String(http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	context.Status(http.StatusOK)
}

func (hdl *Handler) UpdateTask(context *gin.Context) {
	rawID := context.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		context.String(http.StatusBadRequest, "ID must be an integer")
		return
	}

	t := models.TaskUpdateInput{}
	err = context.ShouldBindJSON(&t)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	err = hdl.tasksdb.UpdateTask(id, t)
	if errors.Is(err, sql.ErrNoRows) {
		context.String(http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, t)
}

func (hdl *Handler) CreateTask(context *gin.Context) {
	t := models.TaskCreateInput{}
	err := context.ShouldBindJSON(&t)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	err = hdl.tasksdb.CreateTask(t)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, t)
}
