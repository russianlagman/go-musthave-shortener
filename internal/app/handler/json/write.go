package json

import (
	"errors"
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/service/store"
)

type WriteHandlerRequest struct {
	URL string `json:"url"`
}

type WriteHandlerResponse struct {
	Result string `json:"result"`
}

func WriteHandler(s store.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqObj := &WriteHandlerRequest{}
		err := readBody(r, reqObj)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		uid := handler.ReadContextString(r.Context(), handler.ContextKeyUID{})
		shortURL, err := s.WriteURL(reqObj.URL, uid)
		if err != nil {
			var errConflict *store.ConflictError
			if errors.As(err, &errConflict) {
				respObj := &WriteHandlerResponse{Result: errConflict.ExistingURL}
				writeResponse(w, respObj, http.StatusConflict)
				return
			}
			if errors.Is(err, store.ErrBadInput) {
				writeError(w, err, http.StatusBadRequest)
			} else {
				writeError(w, err, http.StatusInternalServerError)
			}
			return
		}

		respObj := &WriteHandlerResponse{Result: shortURL}
		writeResponse(w, respObj, http.StatusCreated)
	}
}
