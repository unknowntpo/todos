package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
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
