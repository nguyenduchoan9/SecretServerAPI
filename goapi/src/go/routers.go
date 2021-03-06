/*
 * Secret Server
 *
 * This is an API of a secret service. You can save your secret by using the API. You can restrict the access of a secret after the certen number of views or after a certen period of time.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Monitor     bool
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	sh := http.StripPrefix("/v1/ui/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/v1/ui/").Handler(sh)
	summaryVec := BuildSummaryVec("http_response_time_milliseconds", "Latency Percentiles in Milliseconds")
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		if route.Monitor {
			handler = WithMonitoring(handler, route, summaryVec)
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/",
		Index,
		true,
	},

	Route{
		"AddSecret",
		strings.ToUpper("Post"),
		"/v1/secret",
		AddSecret,
		true,
	},

	Route{
		"GetSecretByHash",
		strings.ToUpper("Get"),
		"/v1/secret/{hash}",
		GetSecretByHash,
		true,
	},

	Route{
		"Metrics",
		strings.ToUpper("Get"),
		"/metrics",
		Metrics,
		false,
	},
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	p := promhttp.Handler()
	p.ServeHTTP(w, r)
}
