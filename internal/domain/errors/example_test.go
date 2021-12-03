package errors

import (
	"fmt"
)

func _inner() error {
	const op Op = "inner operation"
	return E(op, New("something goes wrong"))
}

func _middle() error {
	const op Op = "middle operation"
	cnt := 3
	err := _inner()
	if err != nil {
		return E(op, Msg("counter value = %d").Format(cnt), err)
	}
	return nil
}

func _outer() error {
	const op Op = "outer operation"
	const email UserEmail = "alice@example.com"
	err := _middle()
	if err != nil {
		return E(email, op, err)
	}
	return nil
}

func Example() {
	err := _outer()
	fmt.Printf("%v\n", err)
	// Output:
	// alice@example.com: outer operation: middle operation: counter value = 3: inner operation: something goes wrong
}

func ExampleMsg_Format() {
	cnt := 3
	e := E(Msg("counter value = %d").Format(cnt), New("something goes wrong"))
	fmt.Printf("%v", e)
	// Output:
	// counter value = 3: something goes wrong
}
