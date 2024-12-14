package main

import (
	"net/http"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/config"
)

func (app *application) enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsOrigin := config.LoadConfig().Cors // loads a map[string]bool of allowed origins

		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = r.Header.Get("Referer")
		}
		
		_, ok := corsOrigin[origin]
		if !ok {
			app.clientError(w, http.StatusForbidden) // respond with 403 Forbidden
			app.errorLog.Printf("Origin not allowed: %s", origin) // log the origin
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin) // set origin in response header if allowed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}