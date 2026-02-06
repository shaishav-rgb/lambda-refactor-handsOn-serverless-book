package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func handler(ctx context.Context) error {
	log.Println("Starting Lambda – loading AWS config")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v\n", err)
		return err
	}

	host := "rds-proxy.proxy-cd8cmm6ks70a.ap-south-1.rds.amazonaws.com"// RDS Proxy endpoint

	port := 5432
	user := "postgres"
	dbname := "postgres"

	log.Printf("Using RDS Proxy endpoint: %s\n", host)
	log.Println("Generating IAM authentication token")

	token, err := auth.BuildAuthToken(
		ctx,
		fmt.Sprintf("%s:%d", host, port),
		cfg.Region,
		user,
		cfg.Credentials,
	)
	if err != nil {
		log.Printf("Failed to generate IAM auth token: %v\n", err)
		return err
	}

	log.Println("IAM authentication token generated successfully")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, token, dbname,
	)

	log.Println("Opening database connection via RDS Proxy")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Failed to open DB connection: %v\n", err)
		return err
	}
	defer db.Close()

	log.Println("Pinging database to verify connection")

	if err := db.PingContext(ctx); err != nil {
		log.Printf("Database ping failed: %v\n", err)
		return err
	}

	log.Println("✅ SUCCESS: Connected to PostgreSQL via RDS Proxy using IAM authentication")

rows, err := db.QueryContext(ctx, "SELECT * FROM users;")
if err != nil {
	log.Printf("Database query failed: %v\n", err)
	return err
}
	// Optional: iterate over rows to log users
for rows.Next() {
	var id int
	var username, email string
	if err := rows.Scan(&id, &username, &email); err != nil {
		log.Printf("Row scan failed: %v\n", err)
		return err
	}
	log.Printf("User: ID=%d, Username=%s, Email=%s", id, username, email)
}

if err := rows.Err(); err != nil {
	log.Printf("Row iteration error: %v\n", err)
	return err
}

log.Println("✅ SUCCESS: Connected to PostgreSQL via RDS Proxy using IAM authentication and fetched users")

	return nil
}

func main() {
	lambda.Start(handler)
}
