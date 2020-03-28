package onfido_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	onfido "github.com/uw-labs/go-onfido"
)

func ExampleApplicant() {
	ctx := context.Background()

	client, err := onfido.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	if client.Token.Prod() {
		panic("onfido token is only for production use")
	}

	applicant, err := client.CreateApplicant(ctx, onfido.Applicant{
		Email:     "rcrowe@example.co.uk",
		FirstName: "Rob",
		LastName:  "Crowe",
		Addresses: []onfido.Address{
			{
				BuildingNumber: "18",
				Street:         "Wind Corner",
				Town:           "Crawley",
				State:          "West Sussex",
				Postcode:       "NW9 5AB",
				Country:        "GBR",
				StartDate:      "2018-02-10",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	applicant, err = client.GetApplicant(ctx, applicant.ID)
	if err != nil {
		panic(err)
	}

	if err := client.DeleteApplicant(ctx, applicant.ID); err != nil {
		panic(err)
	}
}

func ExampleCheck() {
	ctx := context.Background()

	client, err := onfido.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	if client.Token.Prod() {
		panic("onfido token is only for production use")
	}

	applicant, err := client.CreateApplicant(ctx, onfido.Applicant{
		Email:     "rcrowe@example.co.uk",
		FirstName: "Rob",
		LastName:  "Crowe",
		Addresses: []onfido.Address{
			{
				BuildingNumber: "18",
				Street:         "Wind Corner",
				Town:           "Crawley",
				State:          "West Sussex",
				Postcode:       "NW9 5AB",
				Country:        "GBR",
				StartDate:      "2018-02-10",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	check, err := client.CreateCheck(ctx, applicant.ID, onfido.CheckRequest{
		Type: onfido.CheckTypeStandard,
		Reports: []*onfido.Report{
			{
				Name: onfido.ReportNameDocument,
			},
			{
				Name:    onfido.ReportNameIdentity,
				Variant: onfido.ReportVariantKYC,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Form: %+v\n", check.FormURI)
}

func ExampleError() {
	ctx := context.Background()

	client := onfido.NewClient("")

	err := client.DeleteApplicant(ctx, "123")
	onfidoErr, ok := err.(*onfido.Error)
	if ok {
		fmt.Printf("got error from onfido api: %s\n", onfidoErr)
	}
}

func ExampleClient_NewSdkToken() {
	ctx := context.Background()

	client, err := onfido.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	if client.Token.Prod() {
		panic("onfido token is only for production use")
	}

	applicant, err := client.CreateApplicant(ctx, onfido.Applicant{
		Email:     "rcrowe@example.co.uk",
		FirstName: "Rob",
		LastName:  "Crowe",
		Addresses: []onfido.Address{
			{
				BuildingNumber: "18",
				Street:         "Wind Corner",
				Town:           "Crawley",
				State:          "West Sussex",
				Postcode:       "NW9 5AB",
				Country:        "GBR",
				StartDate:      "2018-02-10",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	t, err := client.NewSdkToken(ctx, applicant.ID, "https://*.onfido.com/documentation/*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %v\n", t.Token)

	if err := client.DeleteApplicant(ctx, applicant.ID); err != nil {
		panic(err)
	}
}

func ExampleClient_ListApplicants() {
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

func ExampleClient_UploadDocument() {
	ctx := context.Background()
	applicantID := "3ac9e550-556f-4c67-84f2-0940c85cbe67"

	client, err := onfido.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	doc, err := os.Open("example_id-card.jpg")
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

func ExampleWebhook() {
	wh, err := onfido.NewWebhookFromEnv()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/webhook/onfido", func(w http.ResponseWriter, req *http.Request) {
		whReq, err := wh.ParseFromRequest(req)
		if err != nil {
			if err == onfido.ErrInvalidWebhookSignature {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid signature"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error occurred"))
			return
		}

		fmt.Fprintf(w, "Webhook: %+v\n", whReq)
	})

	http.ListenAndServe(":8080", nil)
}
