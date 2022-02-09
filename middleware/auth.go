package middleware

import (
	"context"
	"net/http"
	"strings"
	"fmt"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware : to verify all authorized operations
func AuthMiddleware(c *gin.Context) {
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	authorizationToken := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	if idToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id token not available"})
		c.Abort()
		return
	}
	//verify token
	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}
	c.Set("UUID", token.UID)
	c.Next()
}

// func CorsMiddleware(c *gin.Context){
// 	if r.Method == http.MethodOptions {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 		w.Header().Set("Access-Control-Max-Age", "3600")
// 		w.WriteHeader(http.StatusNoContent)
// 		return
// 	}
// 	// Set CORS headers for the main request.
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	next.ServeHTTP(w, r)
// }
