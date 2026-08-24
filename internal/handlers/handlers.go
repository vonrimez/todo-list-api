package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/internal/config"
	"github.com/vonrimez/TaskAPI/internal/database"
	"github.com/vonrimez/TaskAPI/internal/models"
)

type Handler struct {
	tasksdb *database.TasksDB
	usersdb *database.UsersDB
	cfg     *config.Config
}

func GetNewHandler(tasksdb *database.TasksDB, userdb *database.UsersDB, cfg *config.Config) *Handler {
	return &Handler{
		tasksdb: tasksdb, usersdb: userdb, cfg: cfg,
	}
}

func getUserID(context *gin.Context) (int, error) {
	rawID, exists := context.Get("ID")
	if !exists {
		return 0, fmt.Errorf("cannot get param for key \"ID\"")
	}
	id, isInt := rawID.(int)
	if !isInt {
		return 0, fmt.Errorf("param for key \"ID\" should be an integer")
	}
	return id, nil
}

func getTaskID(context *gin.Context) (int, error) {
	rawID := context.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return 0, fmt.Errorf("ID must be an integer")
	}
	return id, nil
}

func (hdl *Handler) UserLogin(context *gin.Context) {
	ul := models.UserLoginInput{}
	err := context.ShouldBindJSON(&ul)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	jwt, err := hdl.usersdb.LoginUser(ul, hdl.cfg.JWT_SECRET)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"JWT": jwt, "USER": hdl.usersdb.GetUserInfo(
			&models.User{
				ID:    ul.ID,
				Name:  ul.Name,
				Email: ul.Email,
				Pass:  ul.Pass,
			},
		),
	})
}

func (hdl *Handler) UserRegister(context *gin.Context) {
	ur := models.UserRegisterInput{}
	err := context.ShouldBindJSON(&ur)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	u, err := hdl.usersdb.CreateUser(ur)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	jwt, err := hdl.usersdb.GetJWT(u.ID, hdl.cfg.JWT_SECRET)
	if err != nil {
		fmt.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"JWT": jwt, "USER": hdl.usersdb.GetUserInfo(u),
	})
}

func (hdl *Handler) GetTasks(context *gin.Context) {
	userID, err := getUserID(context)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}
	tasks, err := hdl.tasksdb.GetTasks(userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK, tasks)
}

func (hdl *Handler) GetTaskById(context *gin.Context) {
	taskID, err := getTaskID(context)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}
	userID, err := getUserID(context)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	task, err := hdl.tasksdb.GetTaskById(taskID, userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, task)
}

func (hdl *Handler) DeleteTask(context *gin.Context) {
	taskID, err := getTaskID(context)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}
	userID, err := getUserID(context)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	dt, err := hdl.tasksdb.DeleteTask(taskID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		context.String(http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	context.String(http.StatusOK, hdl.tasksdb.GetTaskInfo(dt))
}

func (hdl *Handler) UpdateTask(context *gin.Context) {
	taskID, err := getTaskID(context)
	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	t := models.TaskUpdateInput{}
	err = context.ShouldBindJSON(&t)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	userID, err := getUserID(context)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	ut, err := hdl.tasksdb.UpdateTask(t, taskID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		context.String(http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.String(http.StatusOK, hdl.tasksdb.GetTaskInfo(ut))
}

func (hdl *Handler) CreateTask(context *gin.Context) {
	tc := models.TaskCreateInput{}
	err := context.ShouldBindJSON(&tc)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	userID, err := getUserID(context)
	if err != nil {
		fmt.Println(1, err)
		context.Status(http.StatusInternalServerError)
		return
	}

	ct, err := hdl.tasksdb.CreateTask(tc, userID)
	if err != nil {
		fmt.Println(2, err)
		context.Status(http.StatusInternalServerError)
		return
	}
	context.String(http.StatusOK, hdl.tasksdb.GetTaskInfo(ct))
}
