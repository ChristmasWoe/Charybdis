package migrateToFirebase

import (
	"charybdis/api"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

func MigrateCategories(db *gorm.DB, firestore *firestore.Client) {
	var cts []api.Category
	ctx := context.Background()
	db.Table("category").Find(&cts)
	for _, v := range cts {
		_, err := firestore.Collection("categories").Doc(v.Id).Set(ctx, map[string]interface{}{
			"name":        v.Name,
			"description": v.Description,
			"parent_id":   v.ParentId,
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
	}
}
