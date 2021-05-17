package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/unknowntpo/todos/internal/data"
	"github.com/unknowntpo/todos/internal/validator"
)

// createTaskHandler creates a new task.
// Request URL: POST /v1/tasks
func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Done    bool   `json:"done"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	task := &data.Task{
		Title:   input.Title,
		Content: input.Content,
		Done:    input.Done,
	}

	// validate the request
	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tasks.Insert(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new task in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%d", task.ID))

	// Write a JSON response with a 201 Created status code, the task data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"task": task}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTaskHandler shows the detail of specific task.
// Request URL: GET /v1/tasks/:id
func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// task is *data.Task
	task, err := app.models.Tasks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the task ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the existing task record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	task, err := app.models.Tasks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Done    bool   `json:"done"`
	}

	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the request body to the appropriate fields of the task
	// record.
	task.Title = input.Title
	task.Content = input.Content
	task.Done = input.Done

	// Validate the updated task record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated task record to our new Update() method.
	err = app.models.Tasks.Update(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
