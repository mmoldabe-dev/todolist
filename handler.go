package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func TodoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoget(db, w, r)
		case http.MethodPost:
			todopost(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func TodoHandlerID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todogetid(db, w, r)
		case http.MethodPut:
			path := strings.TrimPrefix(r.URL.Path, "/todo/")
			parts := strings.Split(path, "/")
			switch {
			case len(parts) == 1:
				todoupdate(db, w, r)
			case len(parts) == 2 && parts[1] == "open":
				TodoToggle(db, w, r, parts[0], false)
			case len(parts) == 2 && parts[1] == "close":
				TodoToggle(db, w, r, parts[0], true)
			default:
				http.Error(w, "Not Found", http.StatusNotFound)
			}

		case http.MethodDelete:
			tododel(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func todoget(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, done FROM TODOLIST")
	if err != nil {
		slog.Error("error during output db")
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	todos := []Task{}
	for rows.Next() {
		var todo Task
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Done); err != nil {
			slog.Error("Error scan")
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		slog.Error("JSON encode error")
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

func todopost(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var t Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		slog.Error("Error")
		http.Error(w, "invalid decoding", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO TODOLIST (Title, Description, DONE)VALUES(?, ?, ?)", t.Title, t.Description, t.Done)
	if err != nil {
		slog.Error("Error")
		http.Error(w, "Error insert", http.StatusBadRequest)
		return

	}

	id, err := result.LastInsertId()
	if err == nil {
		t.Id = int(id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)

}

func tododel(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todo/")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("DELETE FROM TODOLIST WHERE id=?", id)
	if err != nil {
		slog.Error("erorr")
		http.Error(w, "Erorr delete", http.StatusBadRequest)
		return
	}

}
func todoupdate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var t Task

	id := strings.TrimPrefix(r.URL.Path, "/todo/")

	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		slog.Error("Error")
		http.Error(w, "invalid decoding", http.StatusBadRequest)
		return
	}
	_, err = db.Exec(
		"UPDATE TODOLIST SET Title = ?, Description = ?, DONE = ? WHERE ID = ?",
		t.Title, t.Description, t.Done, id,
	)
	if err != nil {
		slog.Error("Error")
		http.Error(w, "error bd ", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)

}
func TodoToggle(db *sql.DB, w http.ResponseWriter, r *http.Request, id string, doneValue bool) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	result, err := db.Exec("UPDATE TODOLIST SET DONE = ? WHERE ID = ?",
		doneValue, id)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	row := db.QueryRow(
		"SELECT id, title, description, done FROM TODOLIST WHERE id = ?",
		id,
	)
	var t Task
	if err := row.Scan(&t.Id, &t.Title, &t.Description, &t.Done); err != nil {
		slog.Error("Error scanning toggled task", slog.String("error", err.Error()))
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)

}

func todogetid(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todo/")

	if id == "" {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * from TODOLIST WHERE id=? ", id)
	var t Task
	err := row.Scan(&t.Id, &t.Title, &t.Description, &t.Done)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
