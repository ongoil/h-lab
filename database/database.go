package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DBConn is DBConn
	DBConn *gorm.DB
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Connect database
	var err error
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connection successfully")
	fmt.Printf("DBConn = %+v\n", DBConn)

	// Migration
	if err := DataBaseMigration(DBConn); err != nil {
		log.Fatal("Database migration failed:", err)
	}
	fmt.Println("Database migration successfully")

	// Get underlying *sql.DB
	sqlDB, err := DBConn.DB()
	if err != nil {
		log.Fatal("Failed to get database object:", err)
	}

	// Connection pool configuration
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	return DBConn
}
