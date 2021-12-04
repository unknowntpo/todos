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

func Example() {
	err := getHandler()
	fmt.Printf("%v\n", err)
	// Note: You can also use "%+v" to show the full stack trace about error
}

func ExampleMsg_Format() {
	cnt := 3
	e := E(Msg("counter value = %d").Format(cnt), New("something goes wrong"))
	fmt.Printf("%v", e)
}
