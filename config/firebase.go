package config

import (
	"context"
	"fmt"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func SetupFirebase() (*auth.Client, *firestore.Client) {
	serviceAccountKeyFilePath, err := filepath.Abs("config/creds.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println(err.Error())
		panic("Firebase load error")
	}
	//Firebase Auth
	auth, err := app.Auth(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		panic("Firebase load error")
	}

	firestore, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		panic("Firestore load error")
	}

	return auth, firestore
}
