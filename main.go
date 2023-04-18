package main

import (
	"database/sql"
	"fmt"
	"log"
	"purebank/api"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"time"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

func main() {

	sqlStore := OpenDB()

	if sqlStore == nil {
		log.Panic("can't connect to postgres!")
	}

	store := db.NewStore(sqlStore)

	server, err := api.NewServer(store)

	if err != nil {
		log.Fatal("Cannot create server", err)
	}
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config err:", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server: ", err)

		return
	}

}

// OpenDB is a open sql database and limit time to connect to sqldb and wait
// some giving time to connect db
func OpenDB() *sql.DB {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config err:", err)
	}

	for {
		connection, err := ConnectDB(config.DBSource)

		if err != nil {
			fmt.Println("Postgres not ready....")
			counts++
		} else {
			fmt.Println("Connected to a database")
			return connection
		}

		if counts > 10 {
			fmt.Errorf("connection time out %w", err)
			return nil
		}
		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

}

// ConnectDB is try to connect sql db using pgx driver
func ConnectDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil

}
