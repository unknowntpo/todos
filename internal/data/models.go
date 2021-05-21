package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models is a wrapper struct which wraps the TaskModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Tasks TaskModel
	Users UserModel
}

// NewModels returns a Models struct containing the initialized TaskModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Tasks: TaskModel{DB: db},
		Users: UserModel{DB: db},
	}
}
