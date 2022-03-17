package migrateToFirebase

import (
	"charybdis/api"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

func MigrateUsers(db *gorm.DB, firestore *firestore.Client) {
	var users []api.User
	ctx := context.Background()
	db.Table("users").Find(&users)
	for _, v := range users {
		_, err := firestore.Collection("users").Doc(v.UID).Set(ctx, map[string]interface{}{
			"name":  v.Name,
			"email": v.Email,
			"role":  v.Role,
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
	}
}
