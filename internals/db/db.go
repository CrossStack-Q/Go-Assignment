package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func New(addr string, maxOpenConn, maxIdleConn int, maxIdleTime string) (*sql.DB, error) {

	db, err := sql.Open("postgres", addr)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	db.SetMaxIdleConns(maxIdleConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
