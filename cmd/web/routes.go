package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/log-portal", app.upsertUserStats)
	mux.HandleFunc("/log-website", app.updateAppStats)
	mux.HandleFunc("/getAggregatedUserStats", app.getAggregatedUserStats)
	mux.HandleFunc("/getAppStats", app.getAppStats)
	mux.HandleFunc("/getStatsDbHealth", app.getStatsDbHealth)

	handler := app.enableCORS(mux)

	return handler
}