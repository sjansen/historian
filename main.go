package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/sjansen/historian/internal/alb"
	"github.com/sjansen/historian/internal/api"
	"github.com/sjansen/historian/internal/demo"
)

func main() {
	if os.Getenv("HISTORIAN_USE_ALB") == "true" {
		if err := alb.ListenAndServe("", &api.Handler{}); err != nil {
			fmt.Printf("error=%q\n", err.Error())
		}
	} else {
		lambda.Start(demo.APIGHandler)
	}
}
