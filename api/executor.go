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

type WorkHour struct {
	OpenHour     int `json:"openHour"`
	OpenMinutes  int `json:"openMinutes"`
	CloseHour    int `json:"closeHour"`
	CloseMinutes int `json:"closeMinutes"`
}

type Executor struct {
	Id               string         `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" `
	Description      string         `json:"description"`
	DescriptionShort string         `json:"description_short"`
	DateCreated      time.Time      `json:"date_created"`
	ExecutorType     string         `json:"executor_type"`
	ICO              string         `json:"ico" gorm:"index:executors,unique"`
	WebsiteUrl       string         `json:"website_url"`
	MainLocation     GeoPoint       `json:"main_location"`
	WorkingRangeInKm int64          `json:"workingRangeInKm"`
	Categories       pq.StringArray `json:"categories" gorm:"type:text[]"`
	WorkHour         WorkHour       `json:"workHour"`
}

type CreateExecutorInput struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	DescriptionShort string         `json:"description_short"`
	DateCreated      time.Time      `json:"date_created"`
	ExecutorType     string         `json:"executor_type"`
	ICO              string         `json:"ico"`
	WebsiteUrl       string         `json:"website_url"`
	MainLocation     GeoPoint       `json:"main_location"`
	WorkingRangeInKm int64          `json:"workingRangeInKm"`
	Categories       pq.StringArray `json:"categories" gorm:"type:text[]"`
	WorkHour         WorkHour       `json:"workHour"`
}

type ExecutorGetInterface struct {
	Id               string     `gorm:"primaryKey"`
	Name             string     `json:"name" `
	Description      string     `json:"description"`
	DescriptionShort string     `json:"description_short"`
	DateCreated      time.Time  `json:"date_created"`
	ExecutorType     string     `json:"executor_type"`
	ICO              string     `json:"ico"`
	WebsiteUrl       string     `json:"website_url"`
	MainLocation     GeoPoint   `json:"main_location"`
	WorkHour         WorkHour   `json:"workHour"`
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

func (w *WorkHour) Scan(value interface{}) error {
	stringified := fmt.Sprintf("%v", value)
	splitted := strings.Split(stringified, " ")
	open := strings.Split(splitted[0], ":")
	close := strings.Split(splitted[1], ":")
	oh, err := strconv.Atoi(open[0])
	if err != nil {
		return errors.New("couldn't read OpenHours")
	}
	om, err := strconv.Atoi(open[1])
	if err != nil {
		return errors.New("couldn't read OpenMinutes")
	}

	ch, err := strconv.Atoi(close[0])
	if err != nil {
		return errors.New("couldn't read CloseHours")
	}
	cm, err := strconv.Atoi(close[1])
	if err != nil {
		return errors.New("couldn't read CloseMinutes")
	}
	w.OpenHour = oh
	w.OpenMinutes = om
	w.CloseHour = ch
	w.CloseMinutes = cm

	return nil
}

func (w WorkHour) Value() (driver.Value, error) {
	return w.String(), nil
}

func (w *WorkHour) String() string {
	return fmt.Sprintf("%s:%s %s:%s", padNumber(w.OpenHour), padNumber(w.OpenMinutes), padNumber(w.CloseHour), padNumber(w.CloseMinutes))
}

func padNumber(n int) string {
	if n >= 0 && n < 10 {
		return fmt.Sprintf("0%d", n)
	} else {
		return fmt.Sprintf("%d", n)
	}
}

//get Listing
//get firms in category

func getCategoriesByIds(c *gin.Context, ids []string) []category {
	fmt.Println("ids", ids)
	res := make([]category, 0)
	if len(ids) == 0 {
		return res
	}
	db := c.MustGet("db").(*gorm.DB)
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
			ICO:              v.ICO,
			WebsiteUrl:       v.WebsiteUrl,
			MainLocation:     v.MainLocation,
			WorkHour:         v.WorkHour,
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
		ICO:              ex.ICO,
		WebsiteUrl:       ex.WebsiteUrl,
		ExecutorType:     ex.ExecutorType,
		MainLocation:     ex.MainLocation,
		WorkHour:         ex.WorkHour,
		WorkingRangeInKm: ex.WorkingRangeInKm,
		Categories:       ex.Categories,
		Id:               uuid.NewV4().String()}
	ret := db.Create(&exec)
	fmt.Println("create ret", ret.Error, ret.Statement)
	c.JSON(http.StatusOK, gin.H{"data": exec})
}

func UpdateExecutor(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if id == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Invalid or empty id"})
		c.Abort()
		return
	}
	var ex CreateExecutorInput

	if err := c.ShouldBindJSON(&ex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var oldEx Executor
	db.Where("id = ?", id).Find(&oldEx)

	oldEx.Name = ex.Name
	oldEx.Description = ex.Description
	oldEx.DescriptionShort = ex.DescriptionShort
	oldEx.ICO = ex.ICO
	oldEx.WebsiteUrl = ex.WebsiteUrl
	oldEx.ExecutorType = ex.ExecutorType
	oldEx.MainLocation = ex.MainLocation
	oldEx.WorkHour = ex.WorkHour
	oldEx.WorkingRangeInKm = ex.WorkingRangeInKm
	oldEx.Categories = ex.Categories

	ret := db.Save(&oldEx)
	if ret.RowsAffected < 1 {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": fmt.Sprintf("row with id=%s cannot be edit because it doesn't exist", id)})
	} else {
		fmt.Println("rows", ret, ret.RowsAffected)
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}

func GetExecutor(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if id == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Invalid or empty id"})
		c.Abort()
		return
	}
	var ex Executor
	db.Where("id = ?", id).Find(&ex)
	c.JSON(http.StatusOK, gin.H{"data": ex})
}

func DeleteExecutor(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if id == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "Invalid or empty id"})
		c.Abort()
		return
	}
	ret := db.Where("id = ?", id).Delete(&Executor{})

	if db.Error != nil {
		c.Header("Content-Type", "application/json")
		fmt.Println("db error", db.Error, db.Error.Error())
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": "db error"})
	} else if ret.RowsAffected < 1 {
		fmt.Println("error exists")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound,
			gin.H{"Error: ": fmt.Sprintf("row with id=%s cannot be deleted because it doesn't exist", id)})
	} else {
		fmt.Println("rows", ret, ret.RowsAffected)
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
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
