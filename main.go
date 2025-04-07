package main

import (
	"log/slog"
	"os"
)

func main() {
	db, err := InitDB("db.sql")
	if err != nil {
		slog.Error("Error open data base", err)
		os.Exit(1)
	}
	defer db.Close()
}
