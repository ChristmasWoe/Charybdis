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
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	DescriptionShort string         `json:"description_short"`
	DateCreated      time.Time      `json:"date_created"`
	ExecutorType     string         `json:"executor_type"`
	MainLocation     GeoPoint       `sql:"type:geometry(Geometry,4326)" json:"main_location"`
	WorkingRangeInKm int64          `json:"workingRangeInKm"`
	Categories       pq.StringArray `json:"categories" gorm:"type:text[]"`
}

type ExecutorGetInterface struct {
	Id               string     `gorm:"primaryKey"`
	Name             string     `json:"name" `
	Description      string     `json:"description"`
	DescriptionShort string     `json:"description_short"`
	DateCreated      time.Time  `json:"date_created"`
	ExecutorType     string     `json:"executor_type"`
	MainLocation     GeoPoint   `json:"main_location"`
	WorkingRangeInKm int64      `json:"workingRangeInKm"`
	Categories       []category `gorm:"-"`
}

type category struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" `
	Description string `json:"description"`
}

func (p *GeoPoint) Scan(value interface{}) error {
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

//get Listing
//get firms in category

func getCategoriesByIds(c *gin.Context, ids []string) []category {
	db := c.MustGet("db").(*gorm.DB)
	res := make([]category, 0)
	db.Table("category").Find(&res, ids)
	return res
}

func GetExecutors(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	exs := make([]Executor, 0)
	exsExt := make([]ExecutorGetInterface, 0)
	result := db.Find(&exs)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	for _, v := range exs {
		exsExt = append(exsExt, ExecutorGetInterface{
			Id:               v.Id,
			Name:             v.Name,
			Description:      v.Description,
			DescriptionShort: v.DescriptionShort,
			DateCreated:      v.DateCreated,
			ExecutorType:     v.ExecutorType,
			MainLocation:     v.MainLocation,
			WorkingRangeInKm: v.WorkingRangeInKm,
			Categories:       getCategoriesByIds(c, v.Categories),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": exsExt})
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
		DateCreated:      time.Now(),
		ExecutorType:     ex.ExecutorType,
		MainLocation:     ex.MainLocation,
		WorkingRangeInKm: ex.WorkingRangeInKm,
		Categories:       ex.Categories,
		Id:               uuid.NewV4().String()}
	db.Create(&exec)
	c.JSON(http.StatusOK, gin.H{"data": exec})
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
