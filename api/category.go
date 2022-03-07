package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// "github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Category struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" `
	Description string `json:"description"`
	ParentId    string `json:"parent_id"`
}

type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentId    string `json:"parent_id" `
}

type CategoryExt struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" `
	Description string `json:"description"`
	// ParentId    string
	SubCollection []CategoryExt `gorm:"-"`
}

func getSubCategories(parent_id string, c *gin.Context) []CategoryExt {
	fmt.Printf("Looking in db with parent_id: %s \n", parent_id)
	db := c.MustGet("db").(*gorm.DB)
	res := make([]CategoryExt, 0)

	// sqlStatement := `SELECT * FROM category WHERE parent_id = $1;`

	db.Table("category").Where("parent_id = ?", parent_id).Find(&res)

	var s []string
	for _, v := range res {
		s = append(s, v.Name, v.Description, v.Id)
	}

	fmt.Printf("%q\n", s)

	for k, v := range res {
		res[k].SubCollection = make([]CategoryExt, 0)
		res[k].SubCollection = getSubCategories(v.Id, c)

		// v.SubCollection = make([]CategoryExt, 0)
		// v.SubCollection = getSubCategories(v.Id,c)
	}
	// //    .(sqlStatement, parent_id)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// for rows.Next() {
	// 	var sct CategoryExt
	// 	var pid string
	// 	sct.SubCollection = make([]CategoryExt, 0)
	// 	rows.Scan(&sct.Name, &sct.Description, &pid, &sct.Id)
	// 	sct.SubCollection = getSubCategories(sct.Id,&c)
	// 	res = append(res, sct)
	// }
	return res
}

func GetCategories(c *gin.Context) {
	cts := getSubCategories("", c)

	c.JSON(http.StatusOK, gin.H{"data": cts})
	// w.Header().Add("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(cts)
	// ctsBytes, _ := json.MarshalIndent(cts, "", "\t")
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(ctsBytes)
}

func CreateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var ct CreateCategoryInput

	if err := c.ShouldBindJSON(&ct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// r.ParseMultipartForm(0)

	// ct.Name = r.FormValue("name")
	// ct.Description = r.FormValue("description")
	// ct.ParentId = r.FormValue("parent_id")
	// ct.Id = uuid.NewV4().String()

	cat := Category{Name: ct.Name, Description: ct.Description, ParentId: ct.ParentId, Id: uuid.NewV4().String()}
	db.Table("category").Create(&cat)
	c.JSON(http.StatusOK, gin.H{"data": cat})

	// if result := db.Create(&ct); result.Error != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	panic(result.Error)
	// 	fmt.Println(result.Error)
	// }

	// sqlStatement := `INSERT INTO category (name, description, parent_id, id) VALUES ($1, $2, $3, $4)`
	// _, err := db.Exec(sqlStatement, ct.Name, ct.Description, ct.ParentId, ct.Id)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	panic(err)
	// }

	// w.WriteHeader(http.StatusOK)
	// ctBytes, _ := json.MarshalIndent(ct, "", "\t")
	// w.Write(ctBytes)
	// defer db.Close()
}

// func CreateCategory(w http.ResponseWriter, r *http.Request) {
// 	db := config.OpenConnection()
// 	var ct Category
// 	r.ParseMultipartForm(0)

// 	ct.Name = r.FormValue("name")
// 	ct.Description = r.FormValue("description")
// 	ct.ParentId = r.FormValue("parent_id")
// 	ct.Id = uuid.NewV4().String()

// 	if result := db.Create(&ct); result.Error != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		panic(result.Error)
// 		fmt.Println(result.Error)
// 	}

// 	// sqlStatement := `INSERT INTO category (name, description, parent_id, id) VALUES ($1, $2, $3, $4)`
// 	// _, err := db.Exec(sqlStatement, ct.Name, ct.Description, ct.ParentId, ct.Id)
// 	// if err != nil {
// 	// 	w.WriteHeader(http.StatusBadRequest)
// 	// 	panic(err)
// 	// }

// 	w.WriteHeader(http.StatusOK)
// 	ctBytes, _ := json.MarshalIndent(ct, "", "\t")
// 	w.Write(ctBytes)
// 	// defer db.Close()
// }

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
