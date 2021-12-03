package errors

import (
	"database/sql"
	"fmt"
)

func getFromDB() error {
	const op Op = "getFromDB"
	return E(op, KindDatabase, sql.ErrNoRows)
}

func businessLogic() error {
	const op Op = "businessLogic"
	cnt := 3
	err := getFromDB()
	if err != nil {
		return E(op, Msg("counter value = %d").Format(cnt), err)
	}
	return nil
}

func getHandler() error {
	const op Op = "getHandler"
	const email UserEmail = "alice@example.com"
	err := businessLogic()
	if err != nil {
		// Handle error here
		switch {
		case KindIs(err, KindDatabase):
			return E(op, email, err)
		default:
			// Do something else
		}
	}
	return nil
}

// You can also use "%+v" to show the full stack trace about error, like:
// alice@example.com: getHandler: businessLogic: counter value = 3: getFromDB: database error: sql: no rows in result set
// github.com/unknowntpo/todos/internal/domain/errors.E
// 	/Users/unknowntpo/repo/unknowntpo/todos/feat-error/internal/domain/errors/errors.go:230
// github.com/unknowntpo/todos/internal/domain/errors.getFromDB
// 	/Users/unknowntpo/repo/unknowntpo/todos/feat-error/internal/domain/errors/example_test.go:10
// github.com/unknowntpo/todos/internal/domain/errors.businessLogic
// 	/Users/unknowntpo/repo/unknowntpo/todos/feat-error/internal/domain/errors/example_test.go:16
// github.com/unknowntpo/todos/internal/domain/errors.getHandler
// 	/Users/unknowntpo/repo/unknowntpo/todos/feat-error/internal/domain/errors/example_test.go:26
// github.com/unknowntpo/todos/internal/domain/errors.Example
// 	/Users/unknowntpo/repo/unknowntpo/todos/feat-error/internal/domain/errors/example_test.go:41
// testing.runExample
// 	/usr/local/Cellar/go/1.16.4/libexec/src/testing/run_example.go:63
// testing.runExamples
// 	/usr/local/Cellar/go/1.16.4/libexec/src/testing/example.go:44
// testing.(*M).Run
// 	/usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1418
// main.main
// 	_testmain.go:55
// runtime.main
// 	/usr/local/Cellar/go/1.16.4/libexec/src/runtime/proc.go:225
// runtime.goexit
// 	/usr/local/Cellar/go/1.16.4/libexec/src/runtime/asm_amd64.s:1371
func Example() {
	err := getHandler()
	fmt.Printf("%v\n", err)
	// Output:
	// alice@example.com: getHandler: businessLogic: counter value = 3: getFromDB: database error: sql: no rows in result set
}

func ExampleMsg_Format() {
	cnt := 3
	e := E(Msg("counter value = %d").Format(cnt), New("something goes wrong"))
	fmt.Printf("%v", e)
	// Output:
	// counter value = 3: something goes wrong
}
