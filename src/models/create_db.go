package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	"fmt"
	"os"
)

var db *gorm.DB

func init() {
	fmt.Println("init")
	e := godotenv.Load("src/cfg/dev.env")
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn

	db.Debug().AutoMigrate(&Wallet{}, &Transaction{})

	//db.AddForeignKey("from", "wallet(id)", "CASCADE", "CASCADE")
}

func GetDB() *gorm.DB {
	return db
}
