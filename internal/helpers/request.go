package helpers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/unknowntpo/todos/internal/helpers/validator"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func ReadIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// The ReadString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// The ReadCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key count be found, it returns
// the provided default value.
func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// The ReadInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key count be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
// TODO: Use interface to accept different validator.
func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}
