package reactor

import (
	"encoding/json"
	"net/http"
	"sync"
)

var ctxPool = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *Context) WriteJSON(status int, data interface{}) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(status)
	c.w.Write(js)

	return nil
}
