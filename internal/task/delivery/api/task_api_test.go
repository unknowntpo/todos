package api

import (
	"testing"
)

func TestGetAll(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Run("test token", func(t *testing.T) {
		// should return http.StatusOK
		t.Run("valid token", func(t *testing.T) {

		})
		// Should return http.StatusUnauthorized
		t.Run("invalid token", func(t *testing.T) {

		})
	})
	// Should return http.StatusUnauthorized
	t.Run("test params", func(t *testing.T) {
		// /v1/tasks/?title=housework&sort=id&page=1&pagesize=3
		t.Run("test valid params", func(t *testing.T) {

		})
		// /v1/tasks/?title=housework&sort=invalid&id=-33&page=-1&pagesize=-1
		t.Run("test invalid params", func(t *testing.T) {

		})
	})
}

func TestGetByID(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestInsert(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestUpdate(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestDelete(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}
