package main

import (
	"charybdis/api"
	"charybdis/config"
	"charybdis/logging"
	"charybdis/middleware"
	migrations "charybdis/migrations/category"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/autotls"

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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func handleArgs(args []string) {
	mgs := make([]string, 0)
	for i, v := range args {
		if v == "-migrate" {
			mgs = append(mgs, args[i:]...)
		}
	}
	if contains(mgs, "categories") {
		fmt.Println("gonna migrate categories")
		migrations.MigrateCategories()
	}
}

func main() {
	argsWithProg := os.Args[1:]
	handleArgs(argsWithProg)
	// migrations.MigrateCategories()
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
	r.Use(middleware.AuthMiddleware)
	r.Use(logging.Logger())
	// firebaseAuth := config.SetupFirebase()
	// mux := http.NewServeMux()
	//Projects

	r.GET("/user/getAll", api.FindUsers)
	r.GET("/user/get/:email", api.GetUser)
	r.GET("/user/getById/:uid", api.GetUserById)
	r.DELETE("/user/delete/:uid", api.DeleteUser)
	r.POST("/user/create", api.CreateUser)
	r.PUT("/user/edit", api.EditUser)
	r.POST("/category/create", api.CreateCategory)
	r.GET("/category/get", api.GetCategories)
	r.POST("/executor/create", api.CreateExecutor)
	r.GET("/executor/getAll", api.GetExecutors)
	r.PATCH("/executor/update/:id", api.UpdateExecutor)
	r.DELETE("/executor/delete/:id", api.DeleteExecutor)
	r.GET("/executor/get/:id", api.GetExecutor)
	// r.Run(":8080")
	// mux.Handle("/categories/edit", handlerMiddleware(http.HandlerFunc(editProject)))

	// //Labels
	// mux.Handle("/labels/get", handlerMiddleware(http.HandlerFunc(getLabels)))
	// mux.Handle("/labels/create", handlerMiddleware(http.HandlerFunc(createLabel)))
	// mux.Handle("/labels/edit", handlerMiddleware(http.HandlerFunc(editLabel)))
	// // Tasks
	// mux.Handle("/tasks/get", handlerMiddleware(http.HandlerFunc(getTasks)))
	// mux.Handle("/tasks/create", handlerMiddleware(http.HandlerFunc(createTask)))
	// mux.Handle("/tasks/edit", handlerMiddleware(http.HandlerFunc(editTask)))
	// mux.Handle("/task/get", handlerMiddleware(http.HandlerFunc(getTask)))
	// mux.Handle("/task/tick", handlerMiddleware(http.HandlerFunc(tickTask)))
	// mux.Handle("/task/delete", handlerMiddleware(http.HandlerFunc(deleteTask)))
	log.Fatal(autotls.Run(r, "spt-api.xyz", "spt.spt-api.xyz"))
	// r.Run(":8080")

	// err := http.ListenAndServe(":8080", mux)
	// log.Fatal(err)
}
