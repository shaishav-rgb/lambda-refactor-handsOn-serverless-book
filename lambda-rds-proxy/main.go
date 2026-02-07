package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/jackc/pgx/v5/stdlib"
)

 var count int=0

func handler(ctx context.Context) error {
	log.Println("Starting Lambda – loading AWS config")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v\n", err)
		return err
	}

	host := os.Getenv("DB_ENDPOINT") // RDS Proxy endpoint
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	missing := []string{}
	if host == "" {
		missing = append(missing, "DB_HOST")
	}
	if user == "" {
		missing = append(missing, "DB_USER")
	}
	if dbname == "" {
		missing = append(missing, "DB_NAME")
	}
	if port == "" {
		missing = append(missing, "DB_PORT")
	}

	if len(missing) > 0 {
		log.Printf("Missing environment variables: %v\n", missing)
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	log.Printf("Using RDS Proxy endpoint: %s\n", host)
	log.Println("Generating IAM authentication token")

	token, err := auth.BuildAuthToken(
		ctx,
		fmt.Sprintf("%s:%s", host, port),
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
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

	log.Println("✅ SUCCESS: Connected to PostgreSQL via RDS Proxy using IAM authentication and fetched users data from users table")

	return nil
}

func main() {
	// main function gets executed only during initialization of execution environment, later triggers only invoke "handler" function if the execution environment is not yet cleaned up by aws. This is the reason the below log only appears once during initialization of execution environment
	count++;
	log.Printf("lambda execution count is: %v", count)
	lambda.Start(handler)
}
