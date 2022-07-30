package middlewares

import "net/http"

type HttpApiHandler func(response http.ResponseWriter, request *http.Request)

func ApiHandler() HttpApiHandler {
	return func(response http.ResponseWriter, request *http.Request) {

	}
}
