package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ZyoGo/default-ddd-http/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DatabaseConnection(cfg *config.AppConfig) *pgxpool.Pool {
	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig(cfg))
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}

	fmt.Println("Connected to the database!!")

	return connPool
}

func dbConfig(cfg *config.AppConfig) *pgxpool.Config {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Failed to create a config DatabaseConnection, err: ", err)
	}

	dbConfig.MinConns = cfg.Database.MinConnection
	dbConfig.MaxConns = cfg.Database.MaxConnection
	dbConfig.MaxConnLifetime = time.Second * 20
	dbConfig.MaxConnIdleTime = time.Second * 5
	dbConfig.HealthCheckPeriod = time.Second * 30
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	return dbConfig
}
