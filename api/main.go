package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vonrimez/TaskAPI/config"
	"github.com/vonrimez/TaskAPI/internal/database"
	"github.com/vonrimez/TaskAPI/internal/handlers"
	"github.com/vonrimez/TaskAPI/internal/service"
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

	taskRepo := database.GetNewTasksDB(db)
	userRepo := database.GetNewUserDB(db)

	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo, cfg.JWT_SECRET)

	hdl := handlers.GetNewHandler(
		taskService, userService,
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
		tasks := v1.Group("/task")
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
