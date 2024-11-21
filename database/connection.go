package database

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var isBeingConnected bool
var lastPingTime time.Time
var mutex sync.Mutex

func InitDB() *sqlx.DB {
	connect()
	return db
}

func GetDB() *sqlx.DB {
	mutex.Lock()
	defer mutex.Unlock()

	if db == nil {
		connect()
	} else {
		if time.Since(lastPingTime) > time.Minute {
			if err := db.Ping(); err != nil {
				log.Printf("Database connection failed on ping: %s\n", err)
				db.Close()
				connect()
			}
			lastPingTime = time.Now()
		}
	}
	return db
}

func CheckDBConnection() bool {
	if isBeingConnected {
		return false
	}
	mutex.Lock()
	defer mutex.Unlock()
	return tryPingDB()
}

func tryPingDB() bool {
	if db == nil {
		connect()
		return false
	}
	if err := db.Ping(); err != nil {
		log.Printf("Database connection lost: %s. Attempting to reconnect...\n", err)
		db.Close()
		connect()
		return false
	}
	return true
}

func connect() {
	isBeingConnected = true
	defer func() { isBeingConnected = false }()

	var err error
	for retryCount := 0; retryCount < 5; retryCount++ {
		db, err = sqlx.Connect(os.Getenv("DB_TYPE"), os.Getenv("DB_CONN"))
		if err == nil {
			log.Println("Database connection established successfully.")
			db.SetMaxIdleConns(15)
			db.SetMaxOpenConns(70)
			db.SetConnMaxLifetime(time.Minute * 10)
			return
		}
		log.Printf("Failed to connect to database: %s. Retrying in 30 seconds...\n", err)
		time.Sleep(30 * time.Second)
	}
	log.Fatalf("Failed to connect to the database after several attempts: %s\n", err)
}
