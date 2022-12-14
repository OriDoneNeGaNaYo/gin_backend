package infra

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type Database struct {
	DB *gorm.DB
}

func newDatabase() Database {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("MYSQL_ROOT_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_NAME")
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	fmt.Println(URL)
	db, err := gorm.Open(mysql.Open(URL))

	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connection established")
	return Database{
		DB: db,
	}
}

func GetDB() Database {
	LoadEnv()
	return newDatabase()
}
