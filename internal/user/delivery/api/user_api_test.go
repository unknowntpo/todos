package api

import (
	"os"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	t.Skip("TODO: finish the implementation")
	t.Fail()
}
func TestActivateUser(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	t.Skip("TODO: finish the implementation")
	t.Fail()
}
