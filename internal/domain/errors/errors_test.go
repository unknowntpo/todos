package domain

import (
	"database/sql"
	"testing"

	"github.com/pkg/errors"
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
	t.Run("test stack trace", func(t *testing.T) {
		err := outer()
		t.Logf("%+v", err)
		// FIXME: Error.Error should handle showing stack tracing message.
	})
}

func TestError(t *testing.T) {
	t.Run("test duplicated error username and Kind", func(t *testing.T) {

	})
}
func inner() error {
	const op Op = "inner operation"
	// same as doing errors.New("something goes wrong") and wrap it with E()
	return errors.New("something goes wrong")
	//return E(op, "something goes wrong")
}

func middle() error {
	const op Op = "middle operation"
	err := inner()
	if err != nil {
		return E(op, err)
	}
	return nil
}

func outer() error {
	const op Op = "outer operation"
	const user UserName = "alice@example.com"
	err := middle()
	if err != nil {
		return E(user, op, err)
	}
	return nil
}
