package mock

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

type MockMailer struct {
	mock.Mock
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	return MockMailer{sender: sender}
}

func (m MockMailer) Send(recipient, templateFile string, data interface{}) error {
	args := m.Called(recipient, templateFile, data)
	return args.Error(0)
}
