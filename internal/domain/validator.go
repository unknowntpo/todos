package domain

type Validator interface {
	Valid() bool
	AddError(key, message string)
	Check(ok bool, key, message string)
}
