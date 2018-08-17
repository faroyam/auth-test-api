package controller

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/faroyam/auth-test-api/config"
	"github.com/faroyam/auth-test-api/db"
	"github.com/faroyam/auth-test-api/logger"
	"github.com/faroyam/auth-test-api/model"
	"github.com/faroyam/auth-test-api/response"
	"go.uber.org/zap"

	jwt "github.com/dgrijalva/jwt-go"
)

/*
openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub
*/

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

// init reads keys and parses it
func init() {
	logger.ZapLogger.Info("reading keys")
	var err error

	privateKeyBytes, err := ioutil.ReadFile(config.CFG.PrivateRSA)
	if err != nil {
		logger.ZapLogger.Fatal("error while reading private key from file", zap.Error(err))
		return
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		logger.ZapLogger.Fatal("error while generating private key", zap.Error(err))
		return
	}

	publicKeyBytes, err := ioutil.ReadFile(config.CFG.PublicRSA)
	if err != nil {
		logger.ZapLogger.Fatal("error while reading public key from file", zap.Error(err))
		return
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		logger.ZapLogger.Fatal("error while generating public key", zap.Error(err))
		return
	}

}

// Join handler -- registration
func Join(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var user = model.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Request Error"))

		logger.ZapLogger.Info("join request error", zap.Error(err))
		return
	}

	if !user.Validate() {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Invalid Credential"))
		return
	}

	ok, err := db.DB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Internal Error"))

		logger.ZapLogger.Info("db error", zap.Error(err))
		return
	}

	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.NewJSON(response.OK, fmt.Sprintf("Welcome, %s", user.Login)))

		logger.ZapLogger.Info("user added", zap.String("login:", user.Login))
		return
	}
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Login Is Used"))
	return
}

// Auth -- authorization -- gives token to client
func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var user = model.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		logger.ZapLogger.Info("auth request error", zap.Error(err))
		json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Request Error"))
		return
	}

	err = db.DB.CheckAuth(user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Invalid Login or Password"))

		logger.ZapLogger.Info("Auth error", zap.String("login", user.Login), zap.Error(err))
		return
	}

	ipIndeces := strings.LastIndex(r.RemoteAddr, ":")
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["user_name"] = user.Login
	claims["last_user_ip"] = r.RemoteAddr[:ipIndeces]
	tokenString, _ := token.SignedString(privateKey)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewJSON(response.OK, tokenString))

	logger.ZapLogger.Info("Auth successful", zap.String("login", user.Login))
	return
}

// AuthCheck checks if token is valid middleware!
func AuthCheck(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		decryptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		ipIndeces := strings.LastIndex(r.RemoteAddr, ":")
		IP := decryptedToken.Claims.(jwt.MapClaims)["last_user_ip"]

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Authorization Error"))

			logger.ZapLogger.Info("auth error", zap.Error(err))
			return
		}

		if IP != r.RemoteAddr[:ipIndeces] {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Authorization Error"))

			logger.ZapLogger.Info("Auth error", zap.String("reason", "unknown IP"))
			return
		}

		if decryptedToken.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.NewJSON(response.FAILED, "Token Is Not Valid"))

			logger.ZapLogger.Info("Invalid Token")
		}
	})
}

// GetPrivate returns json secret data!
func GetPrivate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewJSON(response.OK, "Private Data"))
	return
}

// GetPublic returns JSON non secret data
func GetPublic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewJSON(response.OK, "Public Data"))
	return
}
