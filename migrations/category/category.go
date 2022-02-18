package migrations

import (
	"charybdis/api"
	"charybdis/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	uuid "github.com/satori/go.uuid"
)

func MigrateCategories() {
	db := config.OpenConnection()

	jf, err := os.Open("migrations/category/categories.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jf.Close()

	byteValue, _ := ioutil.ReadAll(jf)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	for pid, cts := range result {
		for _, v := range cts.([]interface{}) {
			cat := api.Category{Name: v.(string), Description: "", ParentId: pid, Id: uuid.NewV4().String()}
			db.Table("category").Create(&cat)
		}
	}

	fmt.Println("succes")

}
