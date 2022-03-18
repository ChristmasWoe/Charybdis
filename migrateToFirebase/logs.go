package migrateToFirebase

import (
	"charybdis/api"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

func MigrateLogs(db *gorm.DB, firestore *firestore.Client) {
	var logs []api.Log
	ctx := context.Background()
	db.Table("logs").Find(&logs)
	// maxNumber := 10
	// i := 0
	for _, v := range logs {
		if v.Method == "GET" {
			continue
		}
		// i++
		// if i == maxNumber {
		// break
		// }

		_, _, err := firestore.Collection("logs").Add(ctx, map[string]interface{}{
			"uid":        v.Uid,
			"method":     v.Method,
			"controller": v.Controller,
			"action":     v.Action,
			"status":     v.Status,
			"affect_id":  v.AffectId,
			"time":       v.Time,
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
	}
}
