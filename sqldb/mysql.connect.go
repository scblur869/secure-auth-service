package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
//	username = "appuser"
//	password = "1cartman"
//	hostname = "192.168.1.35"
// 	authDb = "authdb"
)

func dsn(dbName string) string {
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	hostname := os.Getenv("MYSQL_HOST")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func Connect2Mysql(dbName string) *sql.DB {
	db, err := sql.Open("mysql", dsn(dbName))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	return db
}
