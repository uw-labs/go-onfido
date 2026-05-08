package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/uw-labs/go-onfido"
)

func main() {
	ctx := context.Background()

	client := onfido.NewClient("")

	err := client.DeleteApplicant(ctx, "123")

	var onfidoErr *onfido.Error
	if errors.As(err, &onfidoErr) {
		fmt.Printf("got error from onfido api: %s\n", onfidoErr)
	}
}
