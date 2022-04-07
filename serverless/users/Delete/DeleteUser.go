package servicedeleteuser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
)

var client *firestore.Client
var authClient *auth.Client

func init() {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: "service-tracker-abfd1"}
	// sa := option.WithCredentialsFile("mamto-a068a-firebase-adminsdk-rxout-88ed79d393.json")
	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	authClient, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

type DeleteUserInput struct {
	Uid string `json:"uid"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := context.Background()
	var dui DeleteUserInput
	err := json.NewDecoder(r.Body).Decode(&dui)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["error"] = "Couldn't retrieve user input"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	// cui.name = r.FormValue("name")
	// cui.email = r.FormValue("email")
	// cui.role = r.FormValue("role")
	// cui.password = r.FormValue("password")

	if dui.Uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["error"] = "Values must be non empty"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	err = authClient.DeleteUser(ctx, dui.Uid)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["error"] = fmt.Sprintf("error deleting user: %v\n", err)
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	_, err = client.Collection("users").Doc(dui.Uid).Delete(ctx)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["error"] = fmt.Sprintf("error deleting user in firestore: %v\n", err)
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Success"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return

}
