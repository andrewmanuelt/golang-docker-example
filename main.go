package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Example struct {
	gorm.Model
	Name    string
	Address string
}

func dbstring() (string, error) {
	// load env
	err := godotenv.Load("./.env")

	if err != nil {
		return "", err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(host.docker.internal)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASS"),
		os.Getenv("MYSQL_DB"),
	)

	return dsn, nil
}

func dbconfig() *gorm.DB {
	dsn, _ := dbstring()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// Migrating table based on struct object
	db.AutoMigrate(&Example{})

	return db
}

func main() {
	// router
	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/data", get_data)
	r.HandleFunc("/save", save_data)

	log.Fatal(http.ListenAndServe(":9000", r))
}

func home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(data)
}

func get_data(w http.ResponseWriter, r *http.Request) {
	// connection
	db := dbconfig()

	var example Example

	db.Find(&example)

	data := map[string]interface{}{
		"status": "ok",
		"data":   &example,
	}

	json.NewEncoder(w).Encode(data)
}

func save_data(w http.ResponseWriter, r *http.Request) {
	// connection
	db := dbconfig()

	var example Example

	db.Find(&example)

	insertdata := Example{
		Name:    "John Doe",
		Address: "902 Rivendell Drive",
	}

	// inserting data
	db.Model(&Example{}).Create(&insertdata)

	data := map[string]interface{}{
		"status": "ok",
		"data":   &example,
	}

	json.NewEncoder(w).Encode(data)
}
