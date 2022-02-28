package api

import (
	"context"
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
	// gorm.Model
	ID        uint      `gorm:"primary_key"`
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

type EditUserInput struct {
	UID  string `json:"uid" binding:"required"`
	Name string `json:"name" binding:"required"`
	Role string `json:"role"  binding:"required"`
}

// FindUsers : Controller for getting all users
func FindUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []User
	db.Table("users").Find(&users)
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
	var u User
	if err := db.Table("users").Where("email = ?", email).First(&u).Error; err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Doesn't match any user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": u})
}

func GetUserById(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Invalid startingIndex on search filter!"})
		c.Abort()
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var u User
	if err := db.Table("users").Where("uid = ?", uid).First(&u).Error; err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Doesn't match any user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": u})
}

func EditUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	client := c.MustGet("firebaseAuth").(*auth.Client)

	var input EditUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := (&auth.UserToUpdate{}).
		DisplayName("John Doe")

	u, err := client.UpdateUser(context.Background(), input.UID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	log.Printf("Successfully updated user: %v\n", u)
	// user := User{Name: input.Name, Email: input.Email, Role: input.Role, UID: uid, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	var cu User

	db.Table("users").Where("uid=?", input.UID).First(&cu)

	db.Table("users").Model(&cu).Updates(map[string]interface{}{"name": input.Name, "role": input.Role, "updated_at": time.Now()})
	c.JSON(http.StatusOK, gin.H{"success": true})
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
	// var user User
	// user.
	user := User{Name: input.Name, Email: input.Email, Role: input.Role, UID: uid, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"data": user})
}
