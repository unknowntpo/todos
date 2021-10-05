package domain

// Mailer is the interface for different mailer implementation.
type Mailer interface {
	Send(recipient, templateFile string, data interface{}) error
}
