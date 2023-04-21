package main

import (
	"database/sql"
	"fmt"
	"github.com/hibiken/asynq"
	logs "github.com/rs/zerolog/log"
	"log"
	"purebank/api"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"purebank/worker"
	"time"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

func main() {

	config, err := util.LoadConfig(".")

	sqlStore := OpenDB(config.DBSource)

	if sqlStore == nil {
		log.Panic("can't connect to postgres!")
	}

	store := db.NewStore(sqlStore)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTasksDestributor(redisOpt)

	server, err := api.NewServer(config, store, taskDistributor)

	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	if err != nil {
		log.Fatal("Cannot load config err:", err)
	}

	go runTaskProcessor(redisOpt, store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server: ", err)

		return
	}

}

func runTaskProcessor(clientOpts asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(clientOpts, store)
	logs.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		logs.Fatal().Err(err).Msg("failed to start task processor")

	}

}

// OpenDB is a open sql database and limit time to connect to sqldb and wait
// some giving time to connect db
func OpenDB(dbsource string) *sql.DB {

	for {
		connection, err := ConnectDB(dbsource)

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
