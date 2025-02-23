package main

import (
	"log"

	"github.com/CrossStack-Q/Go-Assignment/internals/db"
	"github.com/CrossStack-Q/Go-Assignment/internals/store"
)

func main() {

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			addr:         "postgres://postgres:dsa@localhost:5432/sypne?sslmode=disable",
			maxopenConns: 3,
			maxIdleConn:  3,
			maxIdleTime:  "3m",
		},
	}
	db, err := db.New(
		cfg.db.addr,
		int(cfg.db.maxopenConns),
		int(cfg.db.maxopenConns),
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Println("Error in DB Conn", err)
		return
	}

	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.run(app.mount()))
}
