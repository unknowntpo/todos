package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/stretchr/testify/assert"
)

// FIXME: How to deal with foreign key constraint constraint ?
// pq: insert or update on table "tokens" violates foreign key constraint "tokens_user_id_fkey"
// FIXME: How to mock token without touching time.Time ?
func TestInsert(t *testing.T) {
	ctx := context.Background()

	// container and database
	container, db, err := testutil.CreatePostgresTestContainer(ctx, "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(ctx)

	// migration
	mig, err := testutil.NewPgMigrator(db)
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	// test Insert
	wantToken, err := domain.GenerateToken(1, 30*time.Minute, domain.ScopeActivation)
	if err != nil {
		t.Fatal("fail to generate token")
	}

	ctx = context.TODO()

	repo := NewTokenRepo(db)
	err = repo.Insert(ctx, wantToken)
	assert.NoError(t, err)

	// do db query to get inserted token
	var gotToken domain.Token
	query := `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = 1`
	err = db.QueryRowContext(ctx, query).Scan(&gotToken.Hash, &gotToken.UserID, &gotToken.Expiry, &gotToken.Scope)
	assert.NoError(t, err)

	t.Log(gotToken)
	//assert.Equal(t, wantToken, gotToken, "should be equal")
}
