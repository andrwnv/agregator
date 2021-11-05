package handlers

import "net/http"

func CommentRouter(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Hello(w, r)
	})
}

func Hello(w http.ResponseWriter, r *http.Request) {
	payload := map[string]interface{}{
		"data":    "Hello from GO service",
		"success": true,
	}

	respondJson(w, payload, http.StatusOK)
}
