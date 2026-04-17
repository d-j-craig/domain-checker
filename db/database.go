package db

import (
	"context"
	"domain-checker/types"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn
var ctx context.Context

// initates databse for crawler (writing) and server (querying)
func InitDB() error {

	ctx = context.Background()

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

// closed db connection
func CloseDB() {
	Conn.Close(ctx)
}

func InsertDomainData(results types.StatusDB) error{
	// loop through the status db 
	for _, result := range results{
		var domainID int64
		
		// insert data into db domains table
		err := Conn.QueryRow(ctx,
		`INSERT INTO domains (url)
		VALUES ($1)
		ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url
		RETURNING id
		`,result.Name).Scan(&domainID)
		 if err != nil {
            return fmt.Errorf("domain insert failed: %w", err)
        }
		
		// insert data into database status table
		_, err = Conn.Exec(ctx,
		`INSERT INTO status (domain_id, status_code, response_time, response_size, error)
		VALUES ($1, $2, $3, $4, $5)
		`, domainID, result.Status, result.ResponseTime, result.ResponseSize, result.Err)
		if err != nil {
            return fmt.Errorf("status insert failed: %w", err)
        }
	}
	

	return nil

}