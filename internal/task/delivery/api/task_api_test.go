package api

import (
	"testing"
)

func TestGetAll(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Run("SUCCESS", func(t *testing.T) {
		// Should return http.StatusOK, and task

	})
	t.Run("FAIL", func(t *testing.T) {
		// Should return http.StatusUnauthorized
		t.Run("invalid token", func(t *testing.T) {

		})
		// /v1/tasks/?title=housework&sort=invalid&id=-33&page=-1&pagesize=-1
		t.Run("test invalid params", func(t *testing.T) {

		})
	})
}

func TestGetByID(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	// Note: manually set up user in context.
	t.Run("test params - taskID", func(t *testing.T) {
		// Note: Should return http.StatusOK, task
		t.Run("test valid taskID", func(t *testing.T) {

		})
		// Note: Should return http.StatusNotFound
		t.Run("test invalid taskID", func(t *testing.T) {

		})
	})

	t.Fail()
}

func TestInsert(t *testing.T) {
	t.Run("PASS", func(t *testing.T) {

	})
	t.Run("Fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {

		})
		t.Run("failed validation", func(t *testing.T) {

		})
		t.Run("internal server error", func(t *testing.T) {

		})
	})

	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestUpdate(t *testing.T) {
	t.Run("PASS", func(t *testing.T) {

	})
	t.Run("Fail", func(t *testing.T) {
		t.Run("taskid not found", func(t *testing.T) {

		})
		t.Run("failed validation", func(t *testing.T) {
			t.Run("failed validation", func(t *testing.T) {
				t.Run("title", func(t *testing.T) {

				})
				t.Run("content", func(t *testing.T) {

				})
				t.Run("done", func(t *testing.T) {

				})
			})
		})
		t.Run("internal server error", func(t *testing.T) {

		})
	})

	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestDelete(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Run("SUCCESS", func(t *testing.T) {
		// Should return 'task successfully deleted' message
	})
	t.Run("FAIL", func(t *testing.T) {
		// Should receive http.StatusNotFound
		t.Run("invaild taskID param", func(t *testing.T) {

		})
		// taskID is valid, but not found in database.
		t.Run("taskID Not found", func(t *testing.T) {

		})
	})

	t.Fail()
}
