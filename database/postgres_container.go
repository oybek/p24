package database

import (
	"context"
	"database/sql"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Conn      *sql.DB
	Container *testcontainers.Container
}

func RunPostgres(ctx context.Context) (*TestDB, error) {
	username, password, database := "postgres", "postgres", "postgres"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       database,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := postgres.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Host: host,
		User: username,
		Pass: password,
		Name: database,
		Port: port.Port(),
	}
	db, err := Initialize(cfg)
	if err != nil {
		return nil, err
	}

	Migrate(cfg)

	return &TestDB{
		Conn:      db.Conn,
		Container: &postgres,
	}, nil
}

func (testDB *TestDB) Terminate(ctx context.Context) {
	testDB.Conn.Close()
	(*testDB.Container).Terminate(ctx)
}
