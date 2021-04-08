package database

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
)


var Db *pg.DB

func Connect() (con *pg.DB) {
	address := fmt.Sprintf("%s:%s", "localhost", "5432")
	options := &pg.Options{
		Addr:     address,
		User:     "postgres",
		Password: "Kostaq",
		Database: "story",
	}
	Db = pg.Connect(options)
	Db.AddQueryHook(pgdebug.DebugHook{
		Verbose: true,
	})
	err := Db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(" connected ")
	return
}

