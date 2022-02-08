package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// User : Model for User
type User struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique;not null"`
}

// your project directory.
// CreateUserInput : struct for create art post request
type CreateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// FindUsers : Controller for getting all users
func FindUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser : controller for creating new users
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Validate input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create user
	user := User{Name: input.Name, Email: input.Email}
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"data": user})
}
