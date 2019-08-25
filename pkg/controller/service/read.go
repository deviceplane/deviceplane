package service

import (
	"encoding/json"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/validator"
)

func read(r *http.Request, req interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	return validator.Validate(req)
}
