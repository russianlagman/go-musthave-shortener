package basic

import (
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"io/ioutil"
	"net/http"
)

func WriteHandler(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading body: %v", err), http.StatusBadRequest)
			return
		}
		u := string(body)
		uid := handler.ReadContextString(r.Context(), handler.ContextKeyUID{})
		fmt.Printf("\n\n\t\tuid: %v\n\n", uid)
		redirectURL, err := store.WriteURL(u, uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(redirectURL))
	}
}