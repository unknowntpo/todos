package domain

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test building an error
// table-driven test
func TestE(t *testing.T) {
	t.Run("build a record not found error", func(t *testing.T) {
		// define operation
		const op Op = "taskRepo.GetByID"

		// define fakeuserid
		userName := UserName("alice@example.com")
		// define error kind to ErrNotFound
		// assume that we performed an sql query to task database and got sql.ErrNoRows error.
		errFromDB := sql.ErrNoRows

		err := E(op, ErrRecordNotFound, userName, errFromDB)
		assert.Equal(t, "alice@example.com: taskRepo.GetByID: record not found: sql: no rows in result set", err.Error())
	})
	t.Run("build a error with error message only", func(t *testing.T) {
		err := E("some error message")

		assert.Equal(t, "some error message", err.Error())
	})
}
