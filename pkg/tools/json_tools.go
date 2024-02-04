package tools

import (
	"encoding/json"
	"net/http"
)

func Message(status int, description string) map[string]interface{} {
	return map[string]interface{}{"status": status, "description": description}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
