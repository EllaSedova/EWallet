package tools

import (
	"net/http"
)

var NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotFound)
	Respond(w, Message(404, "Этот ресурс не был найден на данном сервере"))
}
