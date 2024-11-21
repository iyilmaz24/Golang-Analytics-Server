package main

import (
	"log"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/env"
)


func main() {
	cfg := config{
		addr:  env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}	

	mux := app.mount()

	log.Fatal(app.run(mux))
}