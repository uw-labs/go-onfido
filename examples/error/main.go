package main

import (
	"context"
	"fmt"

	"github.com/uw-labs/go-onfido"
)

func main() {
	ctx := context.Background()

	client := onfido.NewClient("")

	err := client.DeleteApplicant(ctx, "123")
	onfidoErr, ok := err.(*onfido.Error)
	if ok {
		fmt.Printf("got error from onfido api: %s\n", onfidoErr)
	}
}
