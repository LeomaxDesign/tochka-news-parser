package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4"
)

type PostgresClient struct {
	host     string
	username string
	password string
	dbname   string
	port     int
	DB       *pgx.Conn
}

func New(host, username, password, dbname string, port int) *PostgresClient {
	return &PostgresClient{
		host:     host,
		username: username,
		password: password,
		dbname:   dbname,
		port:     port,
	}
}

func (pc *PostgresClient) Connect() error {
	var err error

	pc.DB, err = pgx.Connect(context.Background(), "host="+pc.host+" user="+pc.username+" dbname="+pc.dbname+" password="+pc.password+" port="+strconv.Itoa(pc.port)+" sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to open sql: %w", err)
	}

	if err = pc.DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping sql: %w", err)
	}

	return nil
}

func (pc *PostgresClient) Disconnect() {
	pc.DB.Close(context.Background())
}
