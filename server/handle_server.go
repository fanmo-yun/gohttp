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

	zap.L().Info("Server Running", zap.String("http", address))
	err := gohttp.ListenAndServe()
	if err != nil {
		zap.L().Fatal("Listen And Serve Fatal", zap.Error(err))
	}
}
