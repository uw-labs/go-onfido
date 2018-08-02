package main

import (
	"context"
	"fmt"

	onfido "github.com/utilitywarehouse/go-onfido"
)

func main() {
	ctx := context.Background()

	client, err := onfido.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	if client.Token.Prod() {
		panic("onfido token is only for production use")
	}

	iter := client.ListApplicants()

	for iter.Next(ctx) {
		fmt.Printf("%+v\n", iter.Applicant())
	}
	if iter.Err() != nil {
		panic(iter.Err())
	}
}
