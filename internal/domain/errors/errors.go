package errors

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	// Op is the operation being performed, usually the name of the method
	// being invoked (userUsecase.Get, tokenRepo.Insert, ...etc.).
	Op Op
	// User is the username of the user attempting the operation.
	User UserName
	// Kind is the class of error, such as Record not found,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// The underlying error that triggered this one, if any.
	Err error
}

// Op describes an operation, usually as the package and method,
// such as "tokenRepo.Insert", or "userUsecase.GetByEmail".
type Op string

// UserName is a string representing a user
type UserName string

// Kind defines the kind of error this is.
type Kind uint8

func (e *Error) Error() string {
	// Build the error message for this level of error.
	sep := ": "

	out := ""
	if e.User != "" {
		out += string(e.User)
		out += sep
	}
	if e.Op != "" {
		out += string(e.Op)
		out += sep
	}

	if e.Kind != 0 {
		out += e.Kind.String()
		out += sep
	}

	if e.Err != nil {
		out += e.Err.Error()
	}

	// FIXME: where does error message from E() goes ?
	// out := e.UserName.String() + sep + e.Op.String() + ... ?
	// recursively call unwrap to unwrap Err.
	// print error message

	return out
}

// Format provide Format method to satisfy fmt.Formatter interface,
// it is used by function like fmt.Printf with verb like: '%+v' to
// display stack trace of a pkg/error Error type.
func (e *Error) Format(f fmt.State, verb rune) {
	if formatter, ok := e.Err.(fmt.Formatter); ok {
		formatter.Format(f, verb)
	} else {
		return
	}
}

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers.
const (
	Other Kind = iota // Unclassified error. This value is not printed in the error message.
	// Maybe moved to httperror.go file ?
	ErrRecordNotFound   // Record not found when we request some resource in database.
	ErrInternal         // Internal server error.
	ErrMethodNotAllowed // Method not allowed error.
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case ErrRecordNotFound:
		return "record not found"
	case ErrInternal:
		return "internal server error"
	case ErrMethodNotAllowed:
		return "method not allowed"
	}
	return "unknown error kind"
}

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	Op
//		The operation of the function who make a call to other function that returns an error.
// 		E.g. taskUsecase.GetAll.
//	UserName
//		The user name of the user attempting the operation.
//	string
//		Treated as an error message.
//	errors.Kind
//		The class of error, such as internal server error, record not found, ...etc.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
//
func E(args ...interface{}) error {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case UserName:
			e.User = arg
		case string:
			e.Err = errors.New(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			e.Err = arg
		case error:
			// if the error implements stackTracer, then it is
			// a pkg/errors error type and does not need to have
			// the stack added
			_, ok := arg.(stackTracer)
			if ok {
				e.Err = arg
			} else {
				e.Err = errors.WithStack(arg)
			}
		default:
			_, file, line, _ := runtime.Caller(1)
			return fmt.Errorf("errors.E: bad call from %s:%d: %v, unknown type %T, value %v in error call", file, line, args, arg, arg)
		}
	}

	return e
}
