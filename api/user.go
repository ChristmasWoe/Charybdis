package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	// "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	// "github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

// User : Model for User
type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// your project directory.
// CreateUserInput : struct for create art post request
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role"  binding:"required"`
	Password string `json:"password" binding:"required"`
}

// FindUsers : Controller for getting all users
func FindUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUser(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Invalid startingIndex on search filter!"})
		c.Abort()
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	// email := c.Query("email")
	fmt.Println(email)
	var u User
	if err := db.Table("users").Where("email = ?", email).First(&u).Error; err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Doesn't match any user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": u})
}

// CreateUser : controller for creating new users
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	client := c.MustGet("firebaseAuth").(*auth.Client)
	// Validate input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := (&auth.UserToCreate{}).
		Email(input.Email).
		EmailVerified(false).
		Password(input.Password).
		DisplayName(input.Name).
		Disabled(false)
	u, err := client.CreateUser(context.Background(), params)
	if err != nil {
		// log.Fatalf("error creating user: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	log.Printf("Successfully created user: %v\n", u)
	uid := u.UserInfo.UID
	// Create user
	user := User{Name: input.Name, Email: input.Email, Role: input.Role, UID: uid, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"data": user})
}
