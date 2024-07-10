package server

import (
	"log"
	"net/http"
)

func ServerRun() {
	gohttp := http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: HandleRouter(),
	}

	log.Println("Listening on http://127.0.0.1:3000...")
	err := gohttp.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
