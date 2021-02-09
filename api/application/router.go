package application

import "net/http"

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc func(res http.ResponseWriter, r *http.Request)
}
