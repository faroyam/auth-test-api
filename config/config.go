package config

import (
	"encoding/json"
	"os"

	"github.com/faroyam/auth-test-api/logger"
	"go.uber.org/zap"
)

// Cfg main config struct
type cfg struct {
	PublicRSA       string `json:"public_rsa"`
	PrivateRSA      string `json:"private_rsa"`
	AppPort         string `json:"app_port"`
	MongoIP         string `json:"mongo_ip"`
	MongoDB         string `json:"mongo_db_name"`
	MongoCollection string `json:"mongo_collection_name"`
}

// CFG main config fle
var CFG = newCFG()

// LoadCFG create Cfg from .json file
func newCFG() cfg {
	var CFG cfg
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		logger.ZapLogger.Fatal("error while loading config.json", zap.Error(err))
	}
	err = json.NewDecoder(configFile).Decode(&CFG)
	if err != nil {
		logger.ZapLogger.Fatal("error while pasring config.json", zap.Error(err))
	}
	return CFG
}
