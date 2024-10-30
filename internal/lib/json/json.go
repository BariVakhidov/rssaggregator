package jsonlib

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

func RespondWithError(w http.ResponseWriter, log *slog.Logger, err error) {
	apiErr, ok := err.(model.APIErr)
	if !ok {
		errResp := model.APIErr{StatusCode: http.StatusInternalServerError, Msg: "internal error"}
		RespondWithJson(w, log, http.StatusInternalServerError, errResp)
	}

	RespondWithJson(w, log, apiErr.StatusCode, apiErr)
}

func RespondWithJson(w http.ResponseWriter, log *slog.Logger, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error("Failed to marshal JSON", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		log.Error("Failed to write JSON response", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
