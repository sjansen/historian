package main

import (
	"fmt"
	"os"

	"github.com/sjansen/historian/internal/alb"
	"github.com/sjansen/historian/internal/api"
	"github.com/sjansen/historian/internal/api/signature"
	"github.com/sjansen/historian/internal/apigw"
	"github.com/sjansen/historian/internal/storage"
)

func main() {
	secret := os.Getenv("HISTORIAN_SECRET")
	table := os.Getenv("HISTORIAN_TABLE")
	useALB := os.Getenv("HISTORIAN_USE_ALB")

	repo, err := storage.NewDynamoDBRepo(table)
	if err != nil {
		fmt.Printf("error=%q\n", err.Error())
		return
	}

	handler := &api.Handler{
		Repo:     repo,
		Verifier: &signature.Verifier{Key: secret},
	}
	if useALB == "true" {
		if err := alb.ListenAndServe("", handler); err != nil {
			fmt.Printf("error=%q\n", err.Error())
		}
	} else {
		if err := apigw.ListenAndServe("", handler); err != nil {
			fmt.Printf("error=%q\n", err.Error())
		}
	}
}
