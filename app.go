package main

import (
	"net/http"

	"github.com/faroyam/auth-test-api/db"
	"github.com/faroyam/auth-test-api/logger"
	"github.com/faroyam/auth-test-api/routes"
	"go.uber.org/zap"
)

func main() {
	logger.ZapLogger.Info("starting server")
	router := routes.NewRouter()
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		logger.ZapLogger.Fatal("error while starting server",
			zap.Error(err),
		)
	}
	defer logger.ZapLogger.Sync()
	defer db.DB.Close()
}
