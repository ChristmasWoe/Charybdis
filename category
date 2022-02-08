package main

import (
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type Category struct {
	Id          string
	Name        string
	Description string
	ParentId    string
}

type CategoryExt struct {
	Id          string
	Name        string
	Description string
	// ParentId    string
	SubCollection []CategoryExt
}

func getSubCategories(parent_id string) []CategoryExt {
	db := OpenConnection()
	res := make([]CategoryExt, 0)

	sqlStatement := `SELECT * FROM category WHERE parent_id = $1;`

	rows, err := db.Query(sqlStatement, parent_id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var sct CategoryExt
		var pid string
		sct.SubCollection = make([]CategoryExt, 0)
		rows.Scan(&sct.Name, &sct.Description, &pid, &sct.Id)
		sct.SubCollection = getSubCategories(sct.Id)
		res = append(res, sct)
	}
	return res
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	cts := getSubCategories("")

	ctsBytes, _ := json.MarshalIndent(cts, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(ctsBytes)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var ct Category
	r.ParseMultipartForm(0)

	ct.Name = r.FormValue("name")
	ct.Description = r.FormValue("description")
	ct.ParentId = r.FormValue("parent_id")
	ct.Id = uuid.NewV4().String()

	sqlStatement := `INSERT INTO category (name, description, parent_id, id) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, ct.Name, ct.Description, ct.ParentId, ct.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	ctBytes, _ := json.MarshalIndent(ct, "", "\t")
	w.Write(ctBytes)
	defer db.Close()
}

// func editProject(w http.ResponseWriter, r *http.Request) {
// 	db := OpenConnection()
// 	r.ParseMultipartForm(0)
// 	var pr Project
// 	pr.Id = r.FormValue("id")
// 	pr.Name = r.FormValue("name")
// 	pr.Description = r.FormValue("description")
// 	pr.Color = r.FormValue("color")
// 	sqlStatement := `UPDATE project SET name = $1, description = $2, color = $3 WHERE id = $4;`
// 	_, err := db.Exec(sqlStatement, pr.Name, pr.Description, pr.Color, pr.Id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		panic(err)
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	prBytes, _ := json.MarshalIndent(pr, "", "\t")
// 	w.Write(prBytes)
// 	defer db.Close()
// }
