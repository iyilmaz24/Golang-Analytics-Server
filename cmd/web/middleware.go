package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/config"
)

func setCorsHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (app *application) enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsOrigin := config.LoadConfig().Cors // loads a map[string]bool of allowed origins

		origin := r.Header.Get("Origin")
		if origin == "" {
			referer := r.Header.Get("Referer")
			if referer != "" {
				if parsedURL, err := url.Parse(referer); err == nil {
					origin = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host) // remove path from referer, only keep scheme and host
				} else {
					app.errorLog.Printf("(cors-middleware) Error parsing referer: %s", err)
				}
			}
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "" {
			if subtle.ConstantTimeCompare([]byte(apiKey), []byte(config.LoadConfig().AdminPassword)) != 1 { // constant-time comparison to prevent timing attacks
				app.errorLog.Printf("(middleware) Invalid API key for request to %s", r.URL.Path)
				app.clientError(w, http.StatusUnauthorized)
				return
			}
			app.infoLog.Printf("API key authenticated request from origin: %s", origin)
			setCorsHeaders(w, origin)
			next.ServeHTTP(w, r) // skip cors origin check if a valid API key is provided
			return
		}

		_, validOrigin := corsOrigin[origin]
		if !validOrigin {
			app.clientError(w, http.StatusForbidden)                                // respond with 403 Forbidden
			app.errorLog.Printf("(cors-middleware) Origin not allowed: %s", origin) // log the origin that was not allowed
			return
		}

		setCorsHeaders(w, origin)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
