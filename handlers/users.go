package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"user-api/models"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with optional filtering by age and sorting by name.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} gin.H
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	db := c.MustGet("db").(*gorm.DB)

	ageStr := c.Query("age")
	var age int
	if ageStr != "" {
		age, _ = strconv.Atoi(ageStr)
	}

	sort := c.Query("sort")
	if sort == "" {
		sort = "name"
	}

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := db.Offset(offset).Limit(limit).Order(sort)

	if age > 0 {
		query = query.Where("age = ?", age)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with unique name and age.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User
	db := c.MustGet("db").(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update an existing user by ID
// @Description Update a user's name and age by their ID. The name must be unique.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User to update"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&user)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by their ID.
// @Tags users
// @Param id path int true "User ID"
// @Success 204 {object} nil
// @Failure 404 {object} gin.H
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
