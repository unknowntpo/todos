package errors

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func inner() error {
	const op Op = "inner operation"
	// same as doing errors.New("something goes wrong") and wrap it with E()
	//return errors.New("something goes wrong")
	return E(op, New("something goes wrong"))
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

func TestOpFormat(t *testing.T) {
	const op Op = "counter.Get - %d"
	var counter int = 3
	out := op.Format(counter)
	assert.Equal(t, "counter.Get - 3", op.Format(counter), "formatted output should be equal")
}

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
		err := E(New("some error message"))

		assert.Equal(t, "some error message", err.Error())
	})
	t.Run("nested error with verb is %s", func(t *testing.T) {
		err := outer()
		want := "alice@example.com: outer operation: middle operation: inner operation: something goes wrong"
		assert.Equal(t, want, fmt.Sprintf("%s", err))
	})
	// TODO: Should we test stacktrace message ?
}

func TestUnwrap(t *testing.T) {
	t.Run("errors.Is", func(t *testing.T) {
		// build nested error
		cause := sql.ErrNoRows
		err := E(Op("repo.Get"), cause)

		// use errors.Is to check the cause is what we want.
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("expect error is %v, but it's not.", cause)
		}
	})
	t.Run("errors.As", func(t *testing.T) {
		// TODO: Test errors.As
	})
}

func TestError(t *testing.T) {
	t.Run("test duplicated error username and Kind", func(t *testing.T) {

	})
}
