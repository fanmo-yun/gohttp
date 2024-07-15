package server

import (
	"context"
	"fmt"
	"gohttp/logger"
	"gohttp/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go func() {
		zap.L().Info("Server Running", zap.String("http", address))
		err := gohttp.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Listen And Serve Fatal", zap.Error(err))
		}
	}()

	ShutdownServer(&gohttp)
}

func ShutdownServer(s *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	zap.L().Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
