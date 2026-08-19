package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vonrimez/TaskAPI/internal/database"
	handlers "github.com/vonrimez/TaskAPI/internal/handling"
)

func main() {

	godotenv.Load()
	DB_URL := os.Getenv("DB_URL")
	PORT := os.Getenv("PORT")

	fmt.Println("Starting server")
	fmt.Println("Connecting with database...")

	db, err := database.EstablishConnection("pgx", DB_URL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Server successfully connected with database")

	tasksdb := database.GetNewTasksDB(db)

	hdl := handlers.GetNewHandler(tasksdb)

	taskApi := gin.New()
	taskApi.Use(gin.Logger())

	v1 := taskApi.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/login", hdl.UserLogin)
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

	err = taskApi.Run(PORT)
	if err != nil {
		panic(err)
	}
}
