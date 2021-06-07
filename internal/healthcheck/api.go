package healthcheck

import (
	"net/http"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
