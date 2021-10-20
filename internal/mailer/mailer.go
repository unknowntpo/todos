package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/domain/errors"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

// Mailer send the message
type Mailer struct {
	config *config.Smtp // config is used to set up the dialer in DefaultSendFunc.
	//letterPaper *Message     // We will fill in the messages we want to send into letterPaper.
}

func New(config *config.Smtp) *Mailer {
	return &Mailer{config}
}

func (m *Mailer) PrepareLetterPaper(recipient, templateName string, data interface{}) (*mail.Message, error) {
	const op errors.Op = "mailer.PrepareLetterPaper"
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateName)
	if err != nil {
		return nil, errors.E(op, err)
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, errors.E(op, err)
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return nil, errors.E(op, err)
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return nil, errors.E(op, err)
	}

	letterPaper := mail.NewMessage()

	letterPaper.SetHeader("To", recipient)
	letterPaper.SetHeader("From", m.config.Sender)
	letterPaper.SetHeader("Subject", subject.String())
	letterPaper.SetBody("text/plain", plainBody.String())
	letterPaper.AddAlternative("text/html", htmlBody.String())

	return letterPaper, nil
}

func (m *Mailer) Send(recipient, templateName string, data interface{}) error {
	const op errors.Op = "mailer.Send"

	lp, err := m.PrepareLetterPaper(recipient, templateName, data)
	if err != nil {
		return errors.E(op, err)
	}

	dialer := mail.NewDialer(m.config.Host, m.config.Port, m.config.Username, m.config.Password)
	dialer.Timeout = 5 * time.Second

	// Try sending the email up to three times before aborting and returning the final
	// error. We sleep for 500 milliseconds between each attempt.
	for i := 1; i <= 3; i++ {
		err = dialer.DialAndSend(lp)
		// If everything worked, return nil.
		if nil == err {
			return nil
		}

		// If it didn't work, sleep for a short time and retry.
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
