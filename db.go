package main

import (
	"database/sql"
)

func InitDB(dataName string) (*sql.DB, error) {
	db, err := sql.Open("github.com/mattn/go-sqlite3", dataName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS TODOLIST(
	   ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Title TEXT NOT NULL,
    Description TEXT,
    DONE BOOLEAN NOT NULL DEFAULT 0
	); 
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
