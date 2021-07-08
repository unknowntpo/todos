package mock

import (
	"testing"
)

func TestGetByID(t *testing.T) {
	t.Run("test id == 1", func(t *testing.T) {

	})

	t.Run("test record not found", func(t *testing.T) {

	})

}

func TestInsert(t *testing.T) {
	// just insert task and assert no errors.

}

func TestUpdate(t *testing.T) {
	t.Run("test successfully update the task", func(t *testing.T) {

	})

	t.Run("test ErrRecordNotFound", func(t *testing.T) {

	})

}

func TestDelete(t *testing.T) {
	t.Run("test successfully delete the task", func(t *testing.T) {

	})

	t.Run("test ErrRecordNotFound", func(t *testing.T) {

	})

}
