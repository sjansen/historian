package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sjansen/historian/internal/demo"
)

func main() {
	lambda.Start(demo.Handler)
}
