package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB   *sql.DB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		var err error
		// Get database path from environment variable or use default
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			// For local development on Windows
			if runtime.GOOS == "windows" {
				dbPath = "./data/app.db"
			} else {
				// For Vercel deployment
				dbPath = "/mnt/data/app.db"
			}
		}
		
		// Ensure directory exists
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			panic("Could not create database directory: " + err.Error())
		}
		
		// Add SQLite parameters
		dbPath = dbPath + "?cache=shared&mode=rwc"
		
		DB, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Print(err)
			panic("Could not connect to database")
		}

		DB.SetMaxOpenConns(1) // Important for in-memory SQLite
		DB.SetMaxIdleConns(1)
		DB.SetConnMaxLifetime(time.Hour)

		_, err = DB.Exec("PRAGMA journal_mode=WAL")
		if err != nil {
			panic("Could not enable WAL mode: " + err.Error())
		}

		_, err = DB.Exec("PRAGMA busy_timeout=5000")
		if err != nil {
			panic("Could not set busy timeout: " + err.Error())
		}

		createTables()
	})
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`

	_, err := DB.Exec(createUsersTable)

	if err != nil {
		fmt.Print(err)
		panic("Could not create users table")
	}

	createEventTable := `
	CREATE TABLE IF NOT EXISTS events(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`

	_, err = DB.Exec(createEventTable)

	if err != nil {
		fmt.Print(err)
		panic("Could not create events table")
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER,
		user_id TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (event_id) REFERENCES events(id)
	)`

	_, err = DB.Exec(createRegistrationTable)

	if err != nil {
		fmt.Print(err)
		panic("Could not create registrations table")
	}

	createAdminsTable := `
	CREATE TABLE IF NOT EXISTS admins(
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`

	_, err = DB.Exec(createAdminsTable)

	if err != nil {
		fmt.Print(err)
		panic("Could not create admins table")
	}
}
