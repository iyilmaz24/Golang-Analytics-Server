package main

import (
	"log"
	"net/http"
	"os"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/config"
	"github.com/iyilmaz24/Go-Analytics-Server/internal/database"
	"github.com/iyilmaz24/Go-Analytics-Server/internal/database/models"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	stats  *models.StatModel
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	appConfig := config.LoadConfig()

	db, err := database.OpenDB(appConfig.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		stats: &models.StatModel{DB: db},
	}

	srv := &http.Server{
		Addr:     appConfig.Port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
        errorLog.Fatal(err)
    }
}