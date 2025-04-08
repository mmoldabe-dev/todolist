package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataName string) (*sql.DB, error) {

	DB, err := sql.Open("sqlite3", dataName)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	err = DB.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS TODOLIST(
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        Title TEXT NOT NULL,
        Description TEXT,
        DONE BOOLEAN NOT NULL DEFAULT 0
    );`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return nil, err
	}

	log.Println("Database initialized successfully.")
	return DB, nil
}
