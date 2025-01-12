package resp

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}
