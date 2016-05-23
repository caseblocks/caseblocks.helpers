package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const StatusUnprocessableEntity = 422
const HTTP_CONTENT_TYPE = "Content-Type"
const HTTP_CONTENT_TYPE_JSON = "application/json"

func RenderJSON(w http.ResponseWriter, payload interface{}) {
	if bytes, err := json.Marshal(payload); err != nil {
		RenderError(w, http.StatusInternalServerError, err)
	} else {
		w.Header().Set(HTTP_CONTENT_TYPE, HTTP_CONTENT_TYPE_JSON)
		w.Write(bytes)
	}
}

func RenderError(w http.ResponseWriter, responseCode int, err error) {
	bytes, _ := json.Marshal(map[string]string{"message": fmt.Sprintf("%s", err)})
	w.WriteHeader(responseCode)
	w.Header().Set(HTTP_CONTENT_TYPE, HTTP_CONTENT_TYPE_JSON)
	w.Write(bytes)
}
