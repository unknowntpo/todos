package domain

type Password interface {
	Set(plaintextPassword string) error
	Matches(plaintextPassword string) (bool, error)
	GetPlainText() string
	GetHash() []byte
}
