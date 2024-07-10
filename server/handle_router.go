package server

import (
	"net/http"
)

func HandleRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			SendHTTPErrorResponse(w, http.StatusMethodNotAllowed)
			return
		}

		if len(r.UserAgent()) == 0 {
			SendHTTPErrorResponse(w, http.StatusForbidden)
			return
		}

		Router(w, r)
	}
}

func Router(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	switch url {
	case "/":
	default:
	}
}
