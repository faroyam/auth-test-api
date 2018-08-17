package main

import (
	"net/http"

	"github.com/faroyam/auth-test-api/db"
	"github.com/faroyam/auth-test-api/logger"
	"github.com/faroyam/auth-test-api/routes"
	"go.uber.org/zap"
)

func main() {
	var err error
	logger.ZapLogger.Info("starting server")
	router := routes.NewRouter()

	err = http.ListenAndServeTLS(":8081", "keys/server.crt", "keys/server.key", router)
	if err != nil {
		logger.ZapLogger.Fatal("error while starting https server", zap.Error(err))
	}
	defer logger.ZapLogger.Sync()
	defer db.DB.Close()
}
