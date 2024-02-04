package app

import (
	t "EWallet/src/tools"
	"net/http"
)

var NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	t.Respond(w, t.Message(false, "This resources was not found on our server"))
}
