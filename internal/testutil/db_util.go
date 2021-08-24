package testutil

import (
	"database/sql"
	"os"
	"path"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

func NewTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("TODOS_DB_DSN"))
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func GetCallerPath() string {
	_, filepath, _, _ := runtime.Caller(1)
	return path.Dir(filepath)
}

// filepath should based on filepath from the project working directory: todos/.
func PrepareSQLQuery(t *testing.T, db *sql.DB, filepath string) func(t *testing.T) {
	t.Helper()
	return func(t *testing.T) {
		// testUtilDir means the file path of testutil package.
		testUtilDir := GetCallerPath()
		wd, _ := os.Getwd()
		t.Log("Getwd", wd)
		t.Log("dir in anonymous func:", testUtilDir)
		query, err := os.ReadFile(testUtilDir + "/../../" + filepath)
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(query))
		if err != nil {
			t.Fatal(err)
		}
	}
}
