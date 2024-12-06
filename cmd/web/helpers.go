package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)


func (app *application) serverError (w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError (w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound (w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For") // check X-Forwarded-For header from NGINX

	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0]) // return the first IP address in the list
	}

	xRealIP := r.Header.Get("X-Real-IP") // check X-Real-IP header from NGINX as fallback
	if xRealIP != "" {
		return xRealIP
	}

	return r.RemoteAddr // lastly, fallback to the remote address
}