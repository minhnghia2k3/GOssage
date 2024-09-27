package database

import (
	"context"
	"database/sql"
	"time"
)

func New(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	idleTime, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(idleTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel() // Release resource if completes before 5 secs

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
