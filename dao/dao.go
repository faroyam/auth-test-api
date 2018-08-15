package dao

import (
	"github.com/faroyam/auth-test-api/logger"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

// DAO Data Access Object
type DAO struct {
	MongoIP            string
	MongoDBName        string
	MongoUserColection string
	Database           string
}

func (d *DAO) connect() (*mgo.Session, error) {
	session, err := mgo.Dial(d.MongoIP)
	return session, err
}

// TODO
func (d *DAO) checkAuth() (bool, error) {
	return false, nil
}

var dao = DAO{}

func init() {
	dao.MongoIP = "192.168.0.11"
	dao.MongoDBName = "test-auth-api"
	dao.MongoUserColection = "users"
	logger.ZapLogger.Info("connecting to db")
	_, err := dao.connect()
	if err != nil {
		logger.ZapLogger.Fatal("connetcing to db error", zap.Error(err))
	}
}
