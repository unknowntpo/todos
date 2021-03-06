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
	const email UserEmail = "alice@example.com"
	err := middle()
	if err != nil {
		return E(email, op, err)
	}
	return nil
}

func TestMsgFormat(t *testing.T) {
	t.Run("formating a message", func(t *testing.T) {
		var counter int = 3
		msg := Msg("current counter value: %d").Format(counter)
		assert.Equal(t, "current counter value: 3", msg.String(), "formatted output should be equal")
	})

	t.Run("formating message in error chain", func(t *testing.T) {
		var counter int = 3
		msg := Msg("current counter value: %d").Format(counter)
		e := E(msg, New("something goes wrong"))
		assert.Equal(t, "current counter value: 3: >> something goes wrong", e.Error(), "error message should be equal")
	})

}

// test building an error
// table-driven test
func TestE(t *testing.T) {
	t.Run("build a KindRecordNotFound error", func(t *testing.T) {
		// define operation
		const op Op = "taskRepo.GetByID"

		// define fake user email
		userEmail := UserEmail("alice@example.com")
		// define error kind to ErrNotFound
		// assume that we performed an sql query to task database and got sql.ErrNoRows error.
		errFromDB := sql.ErrNoRows

		err := E(op, KindRecordNotFound, userEmail, errFromDB)
		assert.Equal(t, "alice@example.com: taskRepo.GetByID: kind record not found: >> sql: no rows in result set", err.Error())
	})
	t.Run("build a error with op and Error.Msg", func(t *testing.T) {
		const op Op = "Counter.Get"
		counter := 3
		err := E(op, Msg("current counter value: %d").Format(counter), New("some error message"))

		assert.Equal(t, "Counter.Get: current counter value: 3: >> some error message", err.Error())
	})

	t.Run("nested error with verb is %s", func(t *testing.T) {
		err := outer()
		want := "alice@example.com: outer operation: >> middle operation: >> inner operation: >> something goes wrong"
		assert.Equal(t, want, fmt.Sprintf("%s", err))
	})
	// TODO: Should we test stacktrace message ?
	t.Run("test if error generated by .E is stackTracer", func(t *testing.T) {
		const op Op = "TestE"
		err := E(op, New("test stacktracer"))

		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		_, ok := err.(stackTracer)
		if !ok {
			t.Fatalf("Error has type %T, and it's not a stackTracer.", err)
		}
	})

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

func TestKindIs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// declare err which is *Error
		err := E(KindInternal)
		// call KindIs()
		isInternal := KindIs(err, KindInternal)
		assert.True(t, isInternal, "kind of error should be KindInternal")
		t.Run("Only inner error has kind specified", func(t *testing.T) {
			inner := E(Op("inner error"), KindFailedValidation, errors.New("something goes wrong"))
			middle := E(Op("middle error"), inner)
			outer := E(Op("outer error"), middle)
			assert.True(t, KindIs(outer, KindFailedValidation), "the outer error with no kind specified should inherit inner error's kind")
		})
		t.Run("No kind specified", func(t *testing.T) {
			inner := E(Op("inner error"), errors.New("something goes wrong"))
			middle := E(Op("middle error"), inner)
			outer := E(Op("outer error"), middle)
			assert.False(t, KindIs(outer, KindFailedValidation), "this should be false because no error kind specified")
		})
	})
	t.Run("Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "want err has the type *Error, got *errors.errorString", r)
			}
		}()

		// call KindIs with error which does not have type *Error.
		_ = KindIs(sql.ErrNoRows, KindInternal)
	})
}
