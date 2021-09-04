package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/stretchr/testify/assert"
)

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

	// Create fake user
	query := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	args := []interface{}{
		"Alice Smith",
		"alice@example.com",
		`\x24326124313224765a682f57676e6246446c4f696757654d62616849756874642f3363526b75576558655469782e374f71524c372e78746570635479`,
		true,
	}

	var user domain.User

	err = db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		t.Fatal(err)
	}

	// test Insert
	wantToken, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
	if err != nil {
		t.Fatal("fail to generate token")
	}

	ctx = context.TODO()

	repo := NewTokenRepo(db)
	err = repo.Insert(ctx, wantToken)
	assert.NoError(t, err)

	// do db query to get inserted token
	var gotToken domain.Token
	query = `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = $1`
	err = db.QueryRowContext(ctx, query, user.ID).Scan(&gotToken.Hash, &gotToken.UserID, &gotToken.Expiry, &gotToken.Scope)
	assert.NoError(t, err)

	// We don't care about Expiry, we just make sure that hash, userid,
	// scope are the same, them we can say that these two token are equal.
	assert.Equal(t, wantToken.Hash, gotToken.Hash, "should be equal")
	assert.Equal(t, wantToken.UserID, gotToken.UserID, "should be equal")
	assert.Equal(t, wantToken.Scope, gotToken.Scope, "should be equal")
}
