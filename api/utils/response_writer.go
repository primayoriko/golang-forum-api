package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponseWriter used to make body's response of the request with(out) JSON format. Could edit status code and header of the response too.
func JSONResponseWriter(w *http.ResponseWriter, statusCode int,
	body interface{}, headerFields map[string]string) error {
	if body != nil {
		(*w).Header().Set("Content-Type", "application/json")
	}

	if headerFields != nil {
		responseHeaders := (*w).Header()
		for k, v := range headerFields {
			responseHeaders.Set(k, v)
		}
	}

	(*w).WriteHeader(statusCode)
	if body != nil {
		jsonString, err := json.Marshal(body)
		if err != nil {
			return err
		}

		(*w).Write(jsonString)
	}

	return nil
}
