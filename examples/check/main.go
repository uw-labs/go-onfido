package main

import (
	"context"
	"fmt"

	"github.com/uw-labs/go-onfido"
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
