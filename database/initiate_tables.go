package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitiateTables(dbPool *pgxpool.Pool) error {
	// Define table creation queries
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			username VARCHAR(30) NOT NULL,
			email VARCHAR(30),
			password VARCHAR(255) NOT NULL,
			role VARCHAR(6) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS unique_email 
			ON users(email, role);
		CREATE UNIQUE INDEX IF NOT EXISTS unique_username 
			ON users(username);
		CREATE INDEX IF NOT EXISTS index_users_id
			ON users (id);
		CREATE INDEX IF NOT EXISTS index_users_name
			ON users USING HASH(lower(username));
		CREATE INDEX IF NOT EXISTS index_users_role
			ON users (role);		
		`,
		// Add more table creation queries here if needed
	}

	// Execute table creation queries
	for _, query := range queries {
		_, err := dbPool.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}
