package main

type task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"discription"`
	Done        bool   `json:"done"`
}
