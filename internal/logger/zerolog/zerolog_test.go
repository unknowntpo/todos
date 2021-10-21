package zerolog

import (
	"bytes"
	"testing"

	"github.com/unknowntpo/todos/internal/domain/errors"

	"github.com/stretchr/testify/assert"
)

func inner() error {
	const op errors.Op = "inner operation"
	return errors.E(op, errors.New("something goes wrong"))
}

func middle() error {
	const op errors.Op = "middle operation"
	err := inner()
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}

func outer() error {
	const op errors.Op = "outer operation"
	const user errors.UserName = "alice@example.com"
	err := middle()
	if err != nil {
		return errors.E(user, op, err)
	}
	return nil
}

func TestPrintInfo(t *testing.T) {
	out := bytes.NewBufferString("")
	log := New(out)
	log.PrintInfo("test PrintInfo", map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	got := out.String()
	assert.Contains(t, got, "time", "log output must contain field 'time'")
	// See https://github.com/sirupsen/logrus/blob/master/logrus_test.go
	// to understand how to make sure our logger has some fields.
	// Output:
	// {"level":"info","key1":"value1","key2":"value2","message":"test PrintInfo"}
}

func TestPrintError(t *testing.T) {
	out := bytes.NewBufferString("")
	log := New(out)

	err := errors.New("something goes wrong")

	log.PrintError(err, map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	got := out.String()
	t.Log(got)
	//assert.Contains(t, got, "time", "log output must contain field 'time'")
	// See https://github.com/sirupsen/logrus/blob/master/logrus_test.go
	// to understand how to make sure our logger has some fields.
	// Output:
	// {"level":"info","key1":"value1","key2":"value2","message":"test PrintInfo"}
}
