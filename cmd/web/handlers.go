package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/database/types"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Welcome to the Golang Analytics Server Catch-All")
	fmt.Fprintln(w, "Use Correct Routes and Methods.")
}

func (app *application) getAggregatedUserStats(w http.ResponseWriter, r *http.Request) {
	// implement the logic to get aggregated user stats
}

func (app *application) getAppStats(w http.ResponseWriter, r *http.Request) {
	// implement the logic to get app stats
}

func (app *application) upsertUserStats(w http.ResponseWriter, r *http.Request) {
	type RequestPayload struct {
		Data   types.UserStat `json:"data"`
	}
	// Sample Payload From TypeScript Frontend: 
	// {
	// 	data: {
	// 	  Ip: "",
	// 	  Location: "",
	// 	  Region: "FL",
	// 	  VD_WebApp: 0,
	// 	  FL_Portal: 1,
	// 	  NM_Portal: 0,
	// 	  TotalVisits: 1,
	// 	  Devices: [
	// 		{
	// 		  Type: "Desktop",
	// 		  OS: "macos",
	// 		  Browser: "chrome",
	// 		},
	// 	  ] as Device[],
	// 	  FirstAccess: new Date().toISOString(),
	// 	  LastAccess: new Date().toISOString(),
	// 	}
	// }

	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		http.Error(w, `{"error": "Failed to parse JSON request body"}`, http.StatusBadRequest)
		return
	}

	clientIp := getClientIP(r)	
	payload.Data.Ip = clientIp

	err = app.stats.UpsertUserStats(&payload.Data)
	if err != nil {
		http.Error(w, `{"error": "Failed to update user stats: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User stats updated successfully"}`))
}

func (app *application) updateAppStats(w http.ResponseWriter, r *http.Request) {
	// implement the logic to update app stats
}

func (app *application) getStatsDbHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := app.stats.CheckHealth()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		app.errorLog.Printf("Could not encode health check response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

