package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDB(ctx context.Context) error {
	// parse the connection string from .env file
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	// disable prepared statements for Supabase pooler
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// connect using the config, use config as we are using the SimpleProtocol (no caching of prepared statments)
	Conn, err = pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}

	// check connection is live
	if err = Conn.Ping(ctx); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func CloseDB(ctx context.Context) {
	Conn.Close(ctx)
}