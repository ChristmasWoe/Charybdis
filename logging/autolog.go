package logging

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Log struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	Uid        string    `json:"uid"`
	Method     string    `json:"method"`
	Controller string    `json:"controller"`
	Action     string    `json:"action"`
	Time       time.Time `json:"time"`
	Latency    int64     `json:"latency"`
	Status     int       `json:"status"`
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var lg Log
		t := time.Now()
		lg.Time = t
		db := c.MustGet("db").(*gorm.DB)

		// Set example variable
		// c.Set("example", "12345")

		// before request

		c.Next()
		// after request
		uid, uidExists := c.Get("UUID")
		if !uidExists {
			return
		}
		lg.Uid = uid.(string)

		latency := time.Since(t)
		lg.Latency = int64(latency / time.Millisecond)

		lg.Method = c.Request.Method
		fp := c.FullPath()
		chunks := strings.Split(fp, "/")
		lg.Controller = chunks[1]
		lg.Action = chunks[2]
		// access the status we are sending
		status := c.Writer.Status()
		lg.Status = status
		db.Table("logs").Create(&lg)
	}
}
