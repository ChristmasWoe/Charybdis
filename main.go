package main

import (
	"charybdis/api"
	"charybdis/config"
	"net/http"

	// "charybdis/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// const (
// 	host     = "172.17.0.2"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "root"
// 	dbname   = "postgres"
// )

func main() {
	r := gin.Default()
	db := config.OpenConnection()
	firebaseAuth := config.SetupFirebase()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("firebaseAuth", firebaseAuth)
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://service-tracker-abfd1.web.app", "http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodPut},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// r.Use(middleware.AuthMiddleware)

	// firebaseAuth := config.SetupFirebase()
	// mux := http.NewServeMux()
	//Projects

	r.GET("/user/getAll", api.FindUsers)
	r.GET("/user/get/:email", api.GetUser)
	r.GET("/user/getById/:uid", api.GetUserById)
	r.POST("/user/create", api.CreateUser)
	r.PUT("/user/edit", api.EditUser)
	r.POST("/category/create", api.CreateCategory)
	r.GET("/category/get", api.GetCategories)
	r.POST("/executor/create", api.CreateExecutor)
	r.GET("/executor/getAll", api.GetExecutors)
	r.PATCH("/executor/update/:id", api.UpdateExecutor)
	r.DELETE("/executor/delete/:id", api.DeleteExecutor)
	r.GET("/executor/get/:id", api.GetExecutor)
	r.Run(":8080")

	// err := http.ListenAndServe(":8080", mux)
	// log.Fatal(err)
}
