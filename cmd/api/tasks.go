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
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Done    bool   `json:done"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}

// showTaskHandler shows the detail of specific task.
// Request URL: GET /v1/tasks/:id
func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
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
		app.serverErrorResponse(w, r, err)
	}
}
