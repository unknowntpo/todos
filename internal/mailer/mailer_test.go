package mailer

import (
	"bytes"
	"testing"
	"time"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {

	/*
		t.Run("send welcome email", func(t *testing.T) {
			// new mailer with config.Config set.
			// create token
			// create data
			//

			m := New("localhost", port, "alice", "pa55word", "TODOs <no-reply@todos.unknowntpo.net>")

			user := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", true)

			token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}

			data := map[string]interface{}{
				"activationToken": token.Plaintext,
				"userID":          user.ID,
			}

			err = m.Send("alice@example.com", "user_welcome.tmpl", data)
			if err != nil {
				t.Fatalf("failed to send welcome email: %v", err)
			}
		})
	*/
}

func TestPrepareLetterPaper(t *testing.T) {
	recipient := "alice@example.com"
	templateName := "user_welcome.tmpl"

	user := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", true)

	token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	data := map[string]interface{}{
		"activationToken": token.Plaintext,
		"userID":          user.ID,
	}

	cfg := &config.Smtp{
		Sender: "TODOs <no-reply@todos.unknowntpo.net>",
	}

	m := New(cfg)

	msg, err := m.PrepareLetterPaper(recipient, templateName, data)
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	msg.WriteTo(buf)

	assert.Containsf(t, buf.String(), cfg.Sender, "sender should be %s", cfg.Sender)
	assert.Containsf(t, buf.String(), "To: alice@example.com", "sender should be %s", cfg.Sender)
	// TODO: Find a proper way to test result of mail.Message!
}
