package middleware

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"scheduling/config"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

// db
var db *sql.DB

// ConnectToDB
func ConnectToDB(dbConfig config.DatabaseConfig) *sql.DB {

	if db != nil {
		return db
	}

	// DSN: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...¶mN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.DBName,
	)
	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error ConnectToDB:", err)
		return nil
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("Error pinging the database: %v", err)
		return nil
	}

	log.Println("Database connection pool initialized successfully.")
	return db
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Println("Error closing the database connection pool:", err)
		} else {
			log.Println("Database connection pool closed.")
		}
	}
}

// LoadConfig reads the TOML configuration file
func LoadConfig(path string) (*config.Config, error) {
	var cfg config.Config
	// Get absolute path for clearer error messages if file not found
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path for %s: %w", path, err)
	}

	log.Printf("Attempting to load configuration from: %s", absPath)

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("error decoding TOML file %s: %w", path, err)
	}
	return &cfg, nil
}
