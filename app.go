package main

import (
	"net/http"

	"github.com/faroyam/auth-test-api/config"
	"github.com/faroyam/auth-test-api/db"
	"github.com/faroyam/auth-test-api/logger"
	"github.com/faroyam/auth-test-api/routes"
	"go.uber.org/zap"
)

func main() {
	var err error
	logger.ZapLogger.Info("starting server")
	router := routes.NewRouter()

	// openssl ecparam -genkey -name secp384r1 -out server.key
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 1

	err = http.ListenAndServeTLS(config.CFG.AppPort, "keys/server.crt", "keys/server.key", router)
	if err != nil {
		logger.ZapLogger.Fatal("error while starting https server", zap.Error(err))
	}
	defer logger.ZapLogger.Sync()
	defer db.DB.Close()
}
