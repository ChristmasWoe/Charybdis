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
	// mux.Handle("/category/get", handlerMiddleware(http.HandlerFunc(getCategories)))
	// mux.Handle("/category/create", handlerMiddleware(http.HandlerFunc(createCategory)))

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

	// log.Fatal(autotls.Run(r, "spt-api.xyz", "spt.spt-api.xyz"))
	r.Run(":8080")

	// err := http.ListenAndServe(":8080", mux)
	// log.Fatal(err)
}

// func handlerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodOptions {
// 			w.Header().Set("Access-Control-Allow-Origin", "*")
// 			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE")
// 			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 			w.Header().Set("Access-Control-Max-Age", "3600")
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		// Set CORS headers for the main request.
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		next.ServeHTTP(w, r)
// 	})
// }

// func OpenConnection() *sql.DB {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return db
// }
