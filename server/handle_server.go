package server

import (
	"fmt"
	"gohttp/logger"
	"gohttp/utils"
	"net/http"

	"go.uber.org/zap"
)

func ServerRun() {
	config := utils.LoadConfig()
	logger.NewLogger(config.Logger)

	address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	gohttp := http.Server{
		Addr:    address,
		Handler: HandleRouter(config),
	}

	zap.L().Info("gohttp: Server running on" + address)
	err := gohttp.ListenAndServe()
	if err != nil {
		zap.L().Sugar().Fatalf("gohttp: Listen And Serve Fatal: %v", err)
	}
}
