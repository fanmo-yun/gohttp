package server

import (
	"errors"
	"gohttp/utils"
	"net/http"
	"os"
)

func HandleRouter() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			SendHTTPErrorResponse(res, http.StatusMethodNotAllowed)
			return
		}

		if len(req.UserAgent()) == 0 {
			SendHTTPErrorResponse(res, http.StatusForbidden)
			return
		}

		Router(res, req)
	}
}

func Router(res http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	switch url {
	case "/":
		if err := utils.VerifyFileExistence("index.html"); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				SendHTTPErrorResponse(res, http.StatusNotFound)
				return
			}
			SendHTTPErrorResponse(res, http.StatusInternalServerError)
			return
		}
		http.ServeFile(res, req, "./html/index.html")
	default:
	}
}
