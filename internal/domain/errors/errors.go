package errors

import (
	"fmt"
	"io"
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
	// UserEmail is the email of the user attempting the operation.
	UserEmail UserEmail
	// Kind is the class of error, such as Record not found,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// Msg is some additional information you want to add.
	Msg Msg
	// The underlying error that triggered this one, if any.
	Err error
}

type Msg string

func (m Msg) String() string {
	return string(m)
}

// Format formats according to a format specifier and return formatted string.
func (m Msg) Format(a ...interface{}) Msg {
	return Msg(fmt.Sprintf(string(m), a...))
}

// Op describes an operation, usually as the package and method,
// such as "tokenRepo.Insert", or "userUsecase.GetByEmail".
type Op string

func (o Op) String() string {
	return string(o)
}

// UserEmail is a string representing a user
type UserEmail string

func (u UserEmail) String() string {
	return string(u)
}

// Kind defines the kind of error this is.
type Kind uint8

func (e *Error) Error() string {
	// Build the error message for this level of error.
	sep := ": "

	out := ""
	if e.UserEmail != "" {
		out += e.UserEmail.String()
		out += sep
	}
	if e.Op != "" {
		out += e.Op.String()
		out += sep
	}

	if e.Kind != 0 {
		out += e.Kind.String()
		out += sep
	}

	if e.Msg != "" {
		out += e.Msg.String()
		out += sep
	}

	if e.Err != nil {
		out += e.Err.Error()
	}

	return out
}

// Format provide Format method to satisfy fmt.Formatter interface,
// it is used by function like fmt.Printf with verb like: '%+v' to
// display stack trace of a pkg/error Error type.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if formatter, ok := e.Err.(fmt.Formatter); ok && s.Flag('+') {
			// Print current level of error message.
			sep := ": "

			out := ""
			if e.UserEmail != "" {
				out += e.UserEmail.String()
				out += sep
			}
			if e.Op != "" {
				out += e.Op.String()
				out += sep
			}

			if e.Kind != 0 {
				out += e.Kind.String()
				out += sep
			}

			if e.Msg != "" {
				out += e.Msg.String()
				out += sep
			}

			io.WriteString(s, out)
			formatter.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *Error) Unwrap() error { return e.Err }

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers.
const (
	KindOther              Kind = iota // Unclassified error. This value is not printed in the error message.
	KindRecordNotFound                 // Record not found when we request some resource in database.
	KindDuplicateEmail                 // Duplicate Email error.
	KindEditConflict                   // Edit conflict while manipulating database.
	KindInvalidCredentials             // Edit conflict while manipulating database.
	KindFailedValidation               //  Failed validation error.
	KindInternal                       // Internal server error.
	KindDatabase                       // Error happened while querying database, this should be treated as subset of internal error and logged it carefully.
)

func (k Kind) String() string {
	switch k {
	case KindOther:
		return "other error"
	case KindRecordNotFound:
		return "record not found"
	case KindDuplicateEmail:
		return "duplicate email"
	case KindEditConflict:
		return "edit conflict"
	case KindInvalidCredentials:
		return "invalid credentials"
	case KindFailedValidation:
		return "failed validation"
	case KindInternal:
		return "internal server error"
	case KindDatabase:
		return "database error"
	}
	return "unknown error kind"
}

// New is the wrapper which calls pkg/errors.New().
func New(text string) error {
	return errors.New(text)
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
//	UserEmail
//		The email of the user attempting the operation.
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
		case UserEmail:
			e.UserEmail = arg
		case Kind:
			//TODO: If Kind not set, maybe we should inherit e.Err's Kind
			e.Kind = arg
		case Msg:
			e.Msg = arg
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

// StackTrace() allow us to print the stacktrace message by calling
// e.Err 's StackTrace() method, if e.Err is nil or e.Err is not a stacktracer,
// we just return nil.
func (e *Error) StackTrace() errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	// No underlying error in e.
	if e.Err == nil {
		return nil
	}
	sterr, ok := e.Err.(stackTracer)
	if !ok {
		return nil
	}
	return sterr.StackTrace()
}

// Is is just a wrapper of pkg/errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is just a wrapper of pkg/errors.As.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// KindIs returns whether err.Kind == kind.
// err must has the type *Error or Error, if not,
// we panic.
// If there are multiple kind specified in error chain,
// we take the outer-most error kind.
func KindIs(err error, kind Kind) bool {
	if e, ok := err.(*Error); ok {
		// fast path: we found a kind in current level of error.
		if e.Kind != KindOther {
			return e.Kind == kind
		}

		// There's no kind specified in current level of error,
		// so we look into e.Err to find out the kind.
		if eInner, ok := e.Err.(*Error); ok && eInner != nil {
			// recursively check if e.Err's kind is equal to kind
			return KindIs(e.Err, kind)
		}

		// there's is no inner error, or inner error is not our custom Error type
		return false
	} else {
		panic(fmt.Sprintf("want err has the type *Error, got %T", err))
	}
}
