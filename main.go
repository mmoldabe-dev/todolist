package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func main() {

	DB, err := InitDB("todo.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer DB.Close()

	// Выводим сообщение, что база данных и таблица успешно созданы
	fmt.Println("Database and table created successfully!")

	http.HandleFunc("/todo", TodoHandler(DB))    // Обработчик для /todo
	http.HandleFunc("/todo/", TodoHandlerID(DB)) // Обработчик для /todo/{id}

	slog.Info("Server starting, Port:8080")
	http.ListenAndServe(":8080", nil) // Запуск сервера
}
