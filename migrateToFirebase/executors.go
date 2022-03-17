package migrateToFirebase

import (
	"charybdis/api"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

func MigrateExecutors(db *gorm.DB, firestore *firestore.Client) {
	var exs []api.Executor
	ctx := context.Background()
	db.Table("executors").Find(&exs)
	// maxNumber := 10
	// i := 0
	for _, v := range exs {
		// i++
		// if i == maxNumber {
		// 	break
		// }

		_, err := firestore.Collection("executors").Doc(v.Id).Set(ctx, map[string]interface{}{
			"name":              v.Name,
			"description":       v.Description,
			"description_short": v.DescriptionShort,
			"executor_type":     v.ExecutorType,
			"ico":               v.ICO,
			"website_url":       v.WebsiteUrl,
			"city":              v.City,
			"address":           v.Address,
			"categories":        v.Categories,
			"workHour":          v.WorkHour,
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
	}
}
