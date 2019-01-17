package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sjansen/historian/internal/demo"
)

func main() {
	if os.Getenv("HISTORIAN_USE_ALB") == "true" {
		lambda.Start(demo.ALBHandler)
	} else {
		lambda.Start(demo.APIGHandler)
	}
}
