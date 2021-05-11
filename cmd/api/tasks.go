package main

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(w, "show the details of task %d\n", id)
}
