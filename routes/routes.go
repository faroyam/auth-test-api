package routes

import (
	"net/http"
	"time"

	"github.com/faroyam/auth-test-api/controller"
	"github.com/faroyam/auth-test-api/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of an API
type Routes []Route

var routes = Routes{
	Route{
		"create_new_user",
		"POST",
		"/v1/join",
		controller.Join,
	},
	Route{
		"token",
		"POST",
		"/v1/auth",
		controller.Auth,
	},
	Route{
		"get_private_data",
		"GET",
		"/v1/private",
		controller.AuthCheck(controller.GetPrivate),
	},
	Route{
		"get_public_data",
		"GET",
		"/v1/public",
		controller.GetPublic,
	},
}

// Logger MIDDLEWARE
func Logger(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.ZapLogger.Info("request",
			zap.String("method", r.Method),
			zap.String("request", r.RequestURI),
			zap.String("client", r.RemoteAddr),
			zap.String("handler_name", name),
			zap.Float64("duration", time.Since(start).Seconds()),
		)
	})
}

// NewRouter configures a router for the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
