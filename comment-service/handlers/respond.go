package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJson(w http.ResponseWriter, payload interface{}, status int) {
	res, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}

		return
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write([]byte(res))

	if err != nil {
		return
	}
}
