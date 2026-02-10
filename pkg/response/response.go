package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
