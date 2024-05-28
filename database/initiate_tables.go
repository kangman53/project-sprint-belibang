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
		`
		CREATE TABLE IF NOT EXISTS merchants (
			id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			name VARCHAR(30) NOT NULL,
			category VARCHAR(100) NOT NULL,
			image_url TEXT NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS index_merchants_id
			ON merchants (id);
		CREATE INDEX IF NOT EXISTS index_merchants_name
			ON merchants USING HASH(lower(name));
		CREATE INDEX IF NOT EXISTS index_merchants_category
			ON merchants (category);
		CREATE INDEX IF NOT EXISTS index_merchants_latitude
			ON merchants (latitude);
		CREATE INDEX IF NOT EXISTS index_merchants_longitude
			ON merchants (longitude);	
		CREATE INDEX IF NOT EXISTS index_merchants_created_at_desc
			ON merchants (created_at DESC);
		CREATE INDEX IF NOT EXISTS index_merchants_created_at_asc
			ON merchants (created_at ASC);
		`,
		`
		CREATE TABLE IF NOT EXISTS items (
			id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			name VARCHAR(30) NOT NULL,
			category VARCHAR(100) NOT NULL,
			price INT NOT NULL,
			image_url TEXT NOT NULL,
			merchant_id VARCHAR(36) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS index_items_merchant_id
			ON items (merchant_id);
		CREATE INDEX IF NOT EXISTS index_items_id
			ON items (id);
		CREATE INDEX IF NOT EXISTS index_items_name
			ON merchants USING HASH(lower(name));
		CREATE INDEX IF NOT EXISTS index_items_category
			ON items (category);
		CREATE INDEX IF NOT EXISTS index_items_created_at_desc
			ON items (created_at DESC);
		CREATE INDEX IF NOT EXISTS index_items_created_at_asc
			ON items (created_at ASC);
		`,
		`
		CREATE TABLE IF NOT EXISTS purchases (
			id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			user_id VARCHAR(36) NOT NULL,
			total_price INT NOT NULL,
			estimated_delivery_time INT NOT NULL,
			distance DOUBLE PRECISION NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			status VARCHAR(7) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			merchant_id VARCHAR(36) NOT NULL,
			is_starting_point BOOLEAN NOT NULL,
			purchase_id VARCHAR(36) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE NO ACTION,
			FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE NO ACTION
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS order_items (
			id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			order_id VARCHAR(36) NOT NULL,
			item_id VARCHAR(36) NOT NULL,
			quantity int NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
		);
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
