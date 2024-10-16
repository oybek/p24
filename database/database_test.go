package database

import (
	"context"
	"os"
	"testing"
	"time"
)

var testdb *TestDB

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := RunPostgres(ctx)
	if err != nil {
		return
	}
	testdb = db

	code := m.Run()
	db.Terminate(ctx)
	os.Exit(code)
}
