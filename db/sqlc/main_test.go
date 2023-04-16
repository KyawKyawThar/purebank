package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"purebank/db/util"
	"testing"
	"time"
)

var testQueries *Queries
var testDB *sql.DB
var counts int64

func TestMain(m *testing.M) {

	testDB = OpenDB()

	if testDB == nil {
		log.Panic("can't connect to postgres!")
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}

// OpenDB is a open sql database and limit time to connect to sqldb and wait
// some giving time to connect db
func OpenDB() *sql.DB {
	config, err := util.LoadConfig("../../")

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
