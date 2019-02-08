package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tumelohq/go-onfido"
)

func main() {
	ctx := context.Background()
	applicantID := "3ac9e550-556f-4c67-84f2-0940c85cbe67"

	client, err := onfido.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	doc, err := os.Open("id-card.jpg")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	document, err := client.UploadDocument(ctx, applicantID, onfido.DocumentRequest{
		File: doc,
		Type: onfido.DocumentTypeIDCard,
		Side: onfido.DocumentSideFront,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", document)
}
