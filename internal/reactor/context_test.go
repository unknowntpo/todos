package reactor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("failed to generate new request: %v", err)
		}

		data := struct {
			S string
			I int
		}{
			"Hello",
			123,
		}
		t.Log(data)

		/*
			data := map[string]string{
				"text":  "134",
				"hello": "Moto",
			}
		*/

		c := &Context{w: rr, r: r}
		err = c.WriteJSON(http.StatusOK, data)
		assert.NoError(t, err)
		t.Log(rr.Result())
		want := `{
	"S": "Hello",
	"I": 123
}
`
		assert.Equal(t, "200 OK", rr.Result().Status)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, want, rr.Body.String())
	})
}
