package json

import (
	"encoding/json"
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"io/ioutil"
	"net/http"
	"net/url"
)

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

type WriteRequest struct {
	URL string `json:"url"`
}

type WriteResponse struct {
	Result string `json:"result"`
}

func WriteHandler(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()

		if err != nil {
			http.Error(w, fmt.Sprintf("error reading body: %v", err), http.StatusBadRequest)
			return
		}

		reqObj := &WriteRequest{}
		err = json.Unmarshal(body, reqObj)
		if err != nil {
			http.Error(w, fmt.Errorf("json read error: %w", err).Error(), http.StatusBadRequest)
			return
		}

		if !IsURL(reqObj.URL) {
			http.Error(w, "bad url", http.StatusBadRequest)
			return
		}
		shortURL, err := store.WriteURL(reqObj.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		resBody, err := json.Marshal(&WriteResponse{Result: shortURL})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = w.Write(resBody)
	}
}
