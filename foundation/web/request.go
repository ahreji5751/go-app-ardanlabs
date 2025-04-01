package web

import (
	"fmt"
	"net/http"

	"github.com/go-json-experiment/json"
)

func Param(r *http.Request, key string) string {
	return r.PathValue(key)
}

type validator interface {
	Validate() error
}

func Decode(r *http.Request, val any) error {
	if err := json.UnmarshalRead(r.Body, val, json.RejectUnknownMembers(false)); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
