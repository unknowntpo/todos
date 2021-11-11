package domain

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortColumn(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	t.Run("test SortColumn with prefix trimmed column", func(t *testing.T) {
		var f Filters

		f.Sort = "-title"
		f.SortSafelist = []string{"title", "-title"}

		got := f.SortColumn()
		want := "title"

		assert.Equal(t, want, got, "should be equal")
	})
	t.Run("test unsafe sorted parameter", func(t *testing.T) {
		assert.Panics(t, func() {
			var f Filters

			f.Sort = "-fiddle"
			f.SortSafelist = []string{"title", "-title"}

			_ = f.SortColumn()
		}, "It should panic")
	})
}

func TestSortDirection(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	t.Run("test ASC", func(t *testing.T) {
		var f Filters

		f.Sort = "title"
		f.SortSafelist = []string{"title", "-title"}

		got := f.SortDirection()
		want := "ASC"

		assert.Equal(t, want, got, "sort direction should be equal")
	})
	t.Run("test DESC", func(t *testing.T) {
		var f Filters

		f.Sort = "-title"
		f.SortSafelist = []string{"title", "-title"}

		got := f.SortDirection()
		want := "DESC"

		assert.Equal(t, want, got, "sort direction should be equal")
	})

}

func TestLimit(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	var f Filters

	f.PageSize = 3

	got := f.Limit()
	want := 3

	assert.Equal(t, want, got, "limit should be equal")
}

func TestOffset(t *testing.T) {
	if os.Getenv("TEST_UNIT") != "1" {
		t.Skip("skipping unit tests")
	}

	var f Filters

	f.PageSize = 3
	f.CurrentPage = 2

	got := f.Limit()
	want := 3

	assert.Equal(t, want, got, "offset should be equal")
}
