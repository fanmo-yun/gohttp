package server

import (
	"fmt"
	"gohttp/utils"
	"log"
	"net/http"
)

func ServerRun() {
	config := utils.LoadConfig()

	address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	gohttp := http.Server{
		Addr:    address,
		Handler: HandleRouter(config),
	}

	log.Println("Listening on http://127.0.0.1:3000...")
	err := gohttp.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
