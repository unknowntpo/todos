package api

import (
	"testing"
)

func TestGetByID(t *testing.T) {
	t.Fail()
}

func TestGetByIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
