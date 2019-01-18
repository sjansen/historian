package main

import (
	"fmt"
	"os"

	"github.com/sjansen/historian/internal/alb"
	"github.com/sjansen/historian/internal/api"
	"github.com/sjansen/historian/internal/apigw"
)

func main() {
	if os.Getenv("HISTORIAN_USE_ALB") == "true" {
		if err := alb.ListenAndServe("", &api.Handler{}); err != nil {
			fmt.Printf("error=%q\n", err.Error())
		}
	} else {
		if err := apigw.ListenAndServe("", &api.Handler{}); err != nil {
			fmt.Printf("error=%q\n", err.Error())
		}
	}
}
