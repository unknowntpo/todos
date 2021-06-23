package api

import (
	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers/validator"
)

// validateTask check if task match the constrains.
func validateTask(v *validator.Validator, task *domain.Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(task.Content != "", "content", "must be provided")
	v.Check(len(task.Content) <= 500, "title", "must not be more than 500 bytes long")
}
