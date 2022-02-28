package logging

import (
	"log"
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
		log.Println("time", lg.Time)
		// after request
		uid, uidExists := c.Get("UUID")
		log.Println("uid exists", uidExists)
		if !uidExists {
			return
		}
		lg.Uid = uid.(string)
		log.Println("uid is", lg.Uid)

		latency := time.Since(t)
		lg.Latency = int64(latency / time.Millisecond)
		log.Println("latency is", lg.Latency)
		log.Print(latency)

		lg.Method = c.Request.Method
		log.Println("Method is", lg.Method)
		fp := c.FullPath()
		log.Println("FullPath is", fp)
		chunks := strings.Split(fp, "/")
		lg.Controller = chunks[1]
		log.Println("Controller is", lg.Controller)
		lg.Action = chunks[2]
		log.Println("Action is", lg.Action)
		// access the status we are sending
		status := c.Writer.Status()
		lg.Status = status
		log.Println("Status is", lg.Status)
		log.Println(status)
		db.Table("logs").Create(&lg)
	}
}
