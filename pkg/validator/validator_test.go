package validator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	t.Run("EmailRX - Match", func(t *testing.T) {
		email := "alice@example.com"

		match := Matches(email, EmailRX)
		assert.True(t, match)
	})

	t.Run("EmailRX - Not Match", func(t *testing.T) {
		email := "aliceasdfsadf"

		match := Matches(email, EmailRX)
		assert.False(t, match)
	})
}

func TestErr(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	v := New()
	v.AddError("key", "value")
	assert.Equal(t, "key: value", v.Err().Error(), "result should be equal")
}
