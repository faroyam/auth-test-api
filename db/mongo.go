package db

import (
	"github.com/faroyam/auth-test-api/config"
	"github.com/faroyam/auth-test-api/logger"
	"github.com/faroyam/auth-test-api/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongo struct {
	MongoIP            string
	MongoDBName        string
	MongoUserColection string
	Database           string
	Session            *mgo.Session
}

func (d *mongo) connect() error {
	var err error
	d.Session, err = mgo.Dial(d.MongoIP)
	return err
}

// Close mgo session
func (d *mongo) Close() {
	d.Session.Close()
}

// CheckAuth return true if user exists & passsword is correct
func (d *mongo) CheckAuth(user model.User) error {
	var exists = model.User{}
	DB.Session.DB(DB.MongoDBName).C(DB.MongoUserColection).Find(bson.M{"login": user.Login}).One(&exists)

	err := bcrypt.CompareHashAndPassword(exists.HashedPassword, []byte(user.Password))
	return err
}

// Join adds new user to db
func (d *mongo) Create(user model.User) (bool, error) {
	var exists = model.User{}
	var err error
	DB.Session.DB(DB.MongoDBName).C(DB.MongoUserColection).Find(bson.M{"login": user.Login}).One(&exists)

	if exists.Login == "" {
		user.ID = bson.NewObjectId()
		user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		err = DB.Session.DB(DB.MongoDBName).C(DB.MongoUserColection).Insert(user)
		return true, err
	}
	return false, nil
}

// DB Data Access Object
var DB = mongo{}

func init() {
	DB.MongoIP = config.CFG.MongoIP
	DB.MongoDBName = config.CFG.MongoDB
	DB.MongoUserColection = config.CFG.MongoCollection
	logger.ZapLogger.Info("connecting to db")
	err := DB.connect()
	if err != nil {
		logger.ZapLogger.Fatal("connetcing to db error", zap.Error(err))
	}
}
