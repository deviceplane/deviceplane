package service

import (
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, ret interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}
