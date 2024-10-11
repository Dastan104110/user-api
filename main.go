package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-api/handlers"
	"user-api/models"
	// Swagger
	"github.com/swaggo/gin-swagger"
	_ "user-api/docs"
)

var db *gorm.DB

// @title User API
// @version 1.0
// @description This is a sample server for managing users.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

func main() {
	var err error
	dsn := "host=localhost user=dastan password=123104110115118 dbname=goproject port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Не удалось выполнить миграцию: %v", err)
	}

	r := gin.Default()
	setupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	v1 := r.Group("/api/v1")
	v1.GET("/users", handlers.GetUsers)
	v1.POST("/users", handlers.CreateUser)
	v1.PUT("/users/:id", handlers.UpdateUser)
	v1.DELETE("/users/:id", handlers.DeleteUser)
}
