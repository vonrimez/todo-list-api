package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/internal/config"
	"github.com/vonrimez/TaskAPI/internal/database"
	"github.com/vonrimez/TaskAPI/internal/handlers"
)

func main() {

	cfg := config.LoadConfig()

	fmt.Println("Starting server")
	fmt.Println("Connecting with database...")

	db, err := database.EstablishConnection("pgx", cfg.DB_URL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Server successfully connected with database")

	tasksdb := database.GetNewTasksDB(db)
	usersdb := database.GetNewUserDB(db)

	hdl := handlers.GetNewHandler(
		tasksdb, usersdb, cfg,
	)

	taskApi := gin.Default()
	taskApi.Use(hdl.StatusLogger)

	v1 := taskApi.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/login", hdl.UserLogin)
			user.POST("/register", hdl.UserRegister)
		}
		tasks := v1.Group("/tasks")
		tasks.Use(hdl.JWTAuth)
		{
			tasks.GET("/list", hdl.GetTasks)
			tasks.GET("/:id", hdl.GetTaskById)
			tasks.DELETE("/:id", hdl.DeleteTask)
			tasks.PUT("/:id", hdl.UpdateTask)
			tasks.POST("/create", hdl.CreateTask)
		}
	}

	fmt.Println("Running server")

	err = taskApi.Run(cfg.PORT)
	if err != nil {
		panic(err)
	}
}
