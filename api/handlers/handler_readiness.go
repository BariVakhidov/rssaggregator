package handlers

import "net/http"

func HandlerReadiness(w http.ResponseWriter, _ *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
