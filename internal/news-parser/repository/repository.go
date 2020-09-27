package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient struct {
	host     string
	username string
	password string
	dbname   string
	port     int
	maxConns int32
	DB       *pgxpool.Pool
}

func New(host, username, password, dbname string, port int, maxConns int32) *PostgresClient {
	return &PostgresClient{
		host:     host,
		username: username,
		password: password,
		dbname:   dbname,
		port:     port,
		maxConns: maxConns,
	}
}

func (pc *PostgresClient) Connect() error {
	var err error

	pc.DB, err = pgxpool.Connect(context.Background(), "host="+pc.host+" user="+pc.username+" dbname="+pc.dbname+" password="+pc.password+" port="+strconv.Itoa(pc.port)+" sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to open sql: %w", err)
	}

	pc.DB.Config().MaxConns = pc.maxConns

	conn, err := pc.DB.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Failed to communcate with database: %s\n", err)
	}
	defer conn.Release()

	if err = conn.Conn().Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping psql: %w", err)
	}

	return nil
}

func (pc *PostgresClient) Disconnect() {
	pc.DB.Close()
}
