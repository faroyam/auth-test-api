package config

type cfg struct {
	PublicRSA       string
	PrivateRSA      string
	AppPort         string
	MongoIP         string
	MongoDB         string
	MongoCollection string
}

// CFG main config fle
var CFG = cfg{
	PublicRSA:       "./keys/app.rsa.pub",
	PrivateRSA:      "./keys/app.rsa",
	AppPort:         ":8081",
	MongoIP:         "192.168.0.11",
	MongoDB:         "test-auth-api",
	MongoCollection: "users",
}
