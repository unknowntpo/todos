package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/data"
)

// createTaskHandler creates a new task.
// Request URL: POST /v1/movies
func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new task")
}

// showTaskHandler shows the detail of specific task.
// Request URL: GET /v1/tasks/:id
func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	task := data.Task{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Do the homework",
		Content:   "p.101 - p.103",
		Done:      false,
		Version:   1,
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
