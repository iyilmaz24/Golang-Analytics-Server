package database

import (
	"database/sql"
	"time"
)

func OpenDB(dsn string) (*sql.DB, error) {

	return nil, nil

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(3) 		// keeps up to 3 idle connections at a time on standby
	db.SetConnMaxLifetime(10 * time.Minute) 		// recycles connections every 10 minutes
	
	return db, nil
}

