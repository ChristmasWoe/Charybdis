package api

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	// "github.com/jinzhu/gorm"
)

type GeoPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Executor struct {
	Id               string         `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" `
	Description      string         `json:"description"`
	DescriptionShort string         `json:"description_short"`
	DateCreated      time.Time      `json:"date_created"`
	ExecutorType     string         `json:"executor_type"`
	MainLocation     GeoPoint       `json:"main_location"`
	WorkingRangeInKm int64          `json:"workingRangeInKm"`
	Categories       pq.StringArray `json:"categories" gorm:"type:text[]"`
}

type CreateExecutorInput struct {
	Name             string         `json:"name"`                                             //`json:"name" binding:"required" `
	Description      string         `json:"description"`                                      //`json:"description" binding:"required"`
	DescriptionShort string         `json:"description_short"`                                //`json:"description_short" binding:"required"`
	DateCreated      time.Time      `json:"date_created"`                                     // `json:"date_created" binding:"required"`
	ExecutorType     string         `json:"executor_type"`                                    // `json:"executor_type" binding:"required"`
	MainLocation     GeoPoint       `sql:"type:geometry(Geometry,4326)" json:"main_location"` // `json:"main_location" binding:"required"`
	WorkingRangeInKm int64          `json:"workingRangeInKm"`                                 // `json:"woringRangeInKm" binding:"required"`
	Categories       pq.StringArray `json:"categories" gorm:"type:text[]"`                    // `json:"categories" binding:"required"`
}

type ExecutorGetInterface struct {
	*Executor
	Categories []category
}

type category struct {
	name        string
	description string
}

func (p *GeoPoint) Scan(value interface{}) error {
	// bytes, ok := value.([]byte)
	// if !ok {
	// 	return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	// }

	stringified := fmt.Sprintf("%v", value)
	fmt.Println("a0", stringified)
	stringified = stringified[6 : len(stringified)-1]
	fmt.Println("a1", stringified)
	splitted := strings.Split(stringified, " ")
	fmt.Println("a3", splitted)
	lng, err := strconv.ParseFloat(splitted[0], 64)
	if err != nil {
		return errors.New(fmt.Sprint("Couldn't read longtitude", value))
	}
	lat, err := strconv.ParseFloat(splitted[1], 64)
	if err != nil {
		return errors.New(fmt.Sprint("Couldn't read latitude", value))
	}
	p.Longitude = lng
	p.Latitude = lat
	return nil

}

func (p GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p *GeoPoint) String() string {
	return fmt.Sprintf("POINT(%v %v)", p.Longitude, p.Latitude)
}

// func getSubCategories(parent_id string, c *gin.Context) []CategoryExt {
// 	fmt.Printf("Looking in db with parent_id: %s \n", parent_id)
// 	db := c.MustGet("db").(*gorm.DB)
// 	res := make([]CategoryExt, 0)

// 	// sqlStatement := `SELECT * FROM category WHERE parent_id = $1;`

// 	db.Table("category").Where("parent_id = ?", parent_id).Find(&res)

// 	var s []string
// 	for _, v := range res {
// 		s = append(s, v.Name, v.Description, v.Id)
// 	}

// 	fmt.Printf("%q\n", s)

// 	for k, v := range res {
// 		res[k].SubCollection = make([]CategoryExt, 0)
// 		res[k].SubCollection = getSubCategories(v.Id, c)

// 		// v.SubCollection = make([]CategoryExt, 0)
// 		// v.SubCollection = getSubCategories(v.Id,c)
// 	}
// 	// //    .(sqlStatement, parent_id)
// 	// // if err != nil {
// 	// // 	log.Fatal(err)
// 	// // }
// 	// for rows.Next() {
// 	// 	var sct CategoryExt
// 	// 	var pid string
// 	// 	sct.SubCollection = make([]CategoryExt, 0)
// 	// 	rows.Scan(&sct.Name, &sct.Description, &pid, &sct.Id)
// 	// 	sct.SubCollection = getSubCategories(sct.Id,&c)
// 	// 	res = append(res, sct)
// 	// }
// 	return res
// }

//get Listing
//get firms in category

func GetExecutors(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	exs := make([]Executor, 0)
	result := db.Find(&exs)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": exs})
}

func CreateExecutor(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var ex CreateExecutorInput

	if err := c.ShouldBindJSON(&ex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Executor", ex.MainLocation)
	exec := Executor{Name: ex.Name,
		Description:      ex.Description,
		DescriptionShort: ex.DescriptionShort,
		DateCreated:      ex.DateCreated,
		ExecutorType:     ex.ExecutorType,
		MainLocation:     ex.MainLocation,
		WorkingRangeInKm: ex.WorkingRangeInKm,
		Categories:       ex.Categories,
		Id:               uuid.NewV4().String()}
	db.Create(&exec)
	c.JSON(http.StatusOK, gin.H{"data": exec})

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
