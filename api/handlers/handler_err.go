package handlers

import "net/http"

func HandlerErr(w http.ResponseWriter, _ *http.Request) {
	respondWithError(w, http.StatusBadRequest, "BAD BAD REQUEST")
}
