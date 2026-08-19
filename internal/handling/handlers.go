package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/internal/database"
)

type Handler struct {
	tasksdb *database.TasksDB
}

func GetNewHandler(tasksdb *database.TasksDB) *Handler {
	return &Handler{tasksdb: tasksdb}
}

func (hdl *Handler) UserLogin(context *gin.Context) {

}

func (hdl *Handler) GetTasks(context *gin.Context) {

}

func (hdl *Handler) GetTaskById(context *gin.Context) {

}

func (hdl *Handler) DeleteTask(context *gin.Context) {

}

func (hdl *Handler) UpdateTask(context *gin.Context) {

}

func (hdl *Handler) CreateTask(context *gin.Context) {

}
