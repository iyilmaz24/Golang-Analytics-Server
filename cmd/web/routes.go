package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/updateUserStats", app.updateUserStats)
	mux.HandleFunc("/updateAppStats", app.updateAppStats)
	mux.HandleFunc("/getAggregatedUserStats", app.getAggregatedUserStats)
	mux.HandleFunc("/getAppStats", app.getAppStats)
	mux.HandleFunc("/getStatsDbHealth", app.getStatsDbHealth)

	return mux
}