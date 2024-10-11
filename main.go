package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-api/handlers"
	"user-api/models"
)

var db *gorm.DB

func main() {
	var err error
	// Строка подключения к базе данных PostgreSQL
	dsn := "host=localhost user=dastan password=123104110115118 dbname=goproject port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Автоматическая миграция
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Не удалось выполнить миграцию: %v", err)
	}

	r := gin.Default()
	setupRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

// setupRoutes - Настройка маршрутов API
func setupRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Set("db", db) // Сохраняем подключение к базе данных в контексте Gin
	})

	v1 := r.Group("/api/v1")
	v1.GET("/users", handlers.GetUsers)          // Получение списка пользователей
	v1.POST("/users", handlers.CreateUser)       // Создание нового пользователя
	v1.PUT("/users/:id", handlers.UpdateUser)    // Обновление пользователя
	v1.DELETE("/users/:id", handlers.DeleteUser) // Удаление пользователя
}
