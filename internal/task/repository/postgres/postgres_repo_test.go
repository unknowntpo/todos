package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set up mock data.
	fakeTime := time.Now()
	row := sqlmock.NewRows([]string{"id", "created_at", "title", "content", "done", "version"}).
		AddRow(1, fakeTime, "task1", "do work", false, 1)
	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = ?").WithArgs(1).WillReturnRows(row)

	repo := NewTaskRepo(db)
	task, err := repo.GetByID(context.TODO(), 1)
	assert.NoError(t, err)

	wantTask := &domain.Task{
		ID:        1,
		CreatedAt: fakeTime,
		Title:     "task1",
		Content:   "do work",
		Done:      false,
		Version:   1,
	}

	assert.Equal(t, wantTask, task)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
