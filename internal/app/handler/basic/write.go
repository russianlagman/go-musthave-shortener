package basic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/service/store"
)

// WriteHandler stores original url and returns the short version
//
//	curl -X POST -d "https://example.org/" http://localhost:8080/
//	http://localhost:8080/xxx
func WriteHandler(s store.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading body: %v", err), http.StatusBadRequest)
			return
		}
		u := string(body)
		uid := handler.ReadContextString(r.Context(), handler.ContextKeyUID{})
		redirectURL, err := s.WriteURL(u, uid)
		if err != nil {
			var errConflict *store.ConflictError
			if errors.As(err, &errConflict) {
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(errConflict.ExistingURL))
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(redirectURL))
	}
}
