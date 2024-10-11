package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"user-api/models"
)

// GetUsers - Получение списка пользователей с фильтрацией и сортировкой
func GetUsers(c *gin.Context) {
	var users []models.User
	db := c.MustGet("db").(*gorm.DB) // Получаем подключение к базе данных

	// Параметры фильтрации
	ageStr := c.Query("age")
	var age int
	if ageStr != "" {
		age, _ = strconv.Atoi(ageStr) // Конвертируем строковый параметр в число
	}

	// Параметры сортировки
	sort := c.Query("sort")
	if sort == "" {
		sort = "name" // По умолчанию сортируем по имени
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10 // Значение по умолчанию
	}
	offset := (page - 1) * limit // Вычисляем смещение

	// Запрос к базе данных с фильтрацией и сортировкой
	query := db.Offset(offset).Limit(limit).Order(sort)

	if age > 0 { // Если передан параметр возраста
		query = query.Where("age = ?", age)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users) // Возвращаем список пользователей
}

// CreateUser - Создание нового пользователя
func CreateUser(c *gin.Context) {
	var user models.User
	db := c.MustGet("db").(*gorm.DB) // Получаем подключение к базе данных

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user) // Возврат созданного пользователя
}

// UpdateUser - Обновление существующего пользователя
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")              // Получение ID пользователя
	db := c.MustGet("db").(*gorm.DB) // Получаем подключение к базе данных

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&user)              // Обновление пользователя
	c.JSON(http.StatusOK, user) // Возврат обновлённого пользователя
}

// DeleteUser - Удаление пользователя
func DeleteUser(c *gin.Context) {
	id := c.Param("id")              // Получение ID пользователя
	db := c.MustGet("db").(*gorm.DB) // Получаем подключение к базе данных

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusNoContent, nil) // Успешное удаление
}
