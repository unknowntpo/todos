package helpers

import (
	"context"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
)

// contextKey represent the contextKey with custom type to avoid naming collision other people's code.
type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey = contextKey("user")

// ContextSetUser returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func ContextSetUser(r *http.Request, user *domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// ContextSetUser retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
func ContextGetUser(r *http.Request) *domain.User {
	user, ok := r.Context().Value(userContextKey).(*domain.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
