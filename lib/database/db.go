package database

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var Container *sqlstore.Container

func Connect() error {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	store, err := sqlstore.New(context.Background(), "sqlite3", "file:majotabi.db?_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL", dbLog)
	if err != nil {
		return err
	}
	Container = store
	fmt.Println("Database Connected (WAL Mode)")
	return nil
}
