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
	USER := getEnv("DB_USER", "root")
	PASS := getEnv("MYSQL_ROOT_PASSWORD", "busdb")
	HOST := getEnv("DB_HOST", "localhost")
	PORT := getEnv("DB_PORT", "3306")
	DBNAME := getEnv("DB_NAME", "bus")
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	fmt.Println(URL)
	db, err := gorm.Open(mysql.Open(URL))

	if err != nil {
		panic("Failed to connect to database!")
	} else {
		fmt.Println("Database connection established")
	}
	return Database{
		DB: db,
	}
}

func GetDB() Database {
	LoadEnv()
	return newDatabase()
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}
