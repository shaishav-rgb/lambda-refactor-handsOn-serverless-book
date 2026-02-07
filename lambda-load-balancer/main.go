package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler() (events.ALBTargetGroupResponse, error) {
	log.Println("Lambda getting invoked hai")
	log.Println("Lambda sending response")
	return events.ALBTargetGroupResponse{
		IsBase64Encoded:   false,
		StatusCode:        200,
		StatusDescription: "200 OK",
		Headers: map[string]string{
			"Set-Cookie":   "cookies",
			"Content-Type": "application/json",
		},
		Body: "Hello from Lambda invoked from ALB",
	}, nil

}

func main() {
	log.Println("lambda execution environment getting created")
	lambda.Start(handler)
}
