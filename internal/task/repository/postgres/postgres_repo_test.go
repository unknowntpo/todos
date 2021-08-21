package postgres

import (
	//"context"

	"context"
	"testing"
	"time"

	//"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/unknowntpo/todos/internal/testutil"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTaskRepo(db)

	rows := sqlmock.NewRows([]string{})

}

// Where to put integration test ?
/*
func TestGetByIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test...")
	}

	db := testutil.NewTestDB(t)

	// TODO: truncate table and inject data manually
	resetTasksTable := `TRUNCATE TABLE tasks;`
	_, err := db.Exec(resetTasksTable)
	assert.NoError(t, err)

	insertTask := `INSERT INTO tasks (title, content, done)
	VALUES ('Do housework', 'The first task!', false);`
	_, err = db.Exec(insertTask)
	assert.NoError(t, err)

	showTasks := `SELECT * FROM tasks;`
	_, err = db.Exec(showTasks)
	assert.NoError(t, err)

	//teardown := testutil.PrepareSQLQuery(t, db, "./migrations/000001_create_tasks_table.down.sql")
	//defer teardown(t)

	// test get id = 1
	repo := NewTaskRepo(db)
	var id int64 = 1
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	task, err := repo.GetByID(ctx, id)
	if assert.NoError(t, err) {
		assert.Equal(t, "Do housework", task.Title, "wrong title")
		assert.Equal(t, "The first task!", task.Content, "wrong content")
		assert.Equal(t, false, task.Done, "wrong state of task: done or not?")

		t.Log(task)
	}
}
*/
