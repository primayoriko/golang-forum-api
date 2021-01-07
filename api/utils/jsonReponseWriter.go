package utils

import (
	"encoding/json"
	"net/http"
)

func JsonResponseWriter(w *http.ResponseWriter, statusCode int, body map[string]interface{}, headerFields map[string]string) error {
	(*w).WriteHeader(statusCode)

	// var m = map[string]interface{}{"complex": 1, "b": []interface{}{"2", "4"}, "c": map[string]interface{}{"int": 1, "string": "giraffe", "bool": true, "float": 4.0}}
	// data, err = json.Marshal(m)
	if body != nil {
		responseHeaders := (*w).Header()
		for k, v := range headerFields {
			responseHeaders.Set(k, v)
		}
	}

	if body != nil {
		jsonString, err := json.Marshal(body)
		if err != nil {
			return err
		}

		(*w).Write(jsonString)
	}

	return nil
}
