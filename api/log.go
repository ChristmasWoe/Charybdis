package api

import (
	"net/http"
	"time"

	// "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	// "github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

// User : Model for User
type Log struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	Uid        string    `json:"uid"`
	Method     string    `json:"method"`
	Controller string    `json:"controller"`
	Action     string    `json:"action"`
	Time       time.Time `json:"time"`
	Latency    int64     `json:"latency"`
	Status     int       `json:"status"`
	AffectId   string    `json:"affect_id"`
}

func GetLogsById(c *gin.Context) {
	uid := c.Param("uid")

	if uid == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error": "Invalid startingIndex on search filter!"})
		c.Abort()
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	lgs := make([]Log, 0)

	if err := db.Table("logs").Where("uid = ?", uid).Where("method IN ?", []string{"POST", "DELETE", "PUT", "PATCH"}).Find(&lgs).Error; err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"Error": "Doesn't match any user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lgs})
}
