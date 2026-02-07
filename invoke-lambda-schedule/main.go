package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() (string, error) {
	log.Println("invoking scheduled lambda function")
	return "Hello", nil
}

func main() {
	lambda.Start(handler)
}
