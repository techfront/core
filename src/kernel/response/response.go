package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
	Meta map[string]interface{} `json:"_meta"`
}

func Send(w http.ResponseWriter, status int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(b)

	return nil
}
