package main

import (
	"context"
	"fmt"


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

	t, err := client.NewSdkToken(ctx, applicant.ID, "https://*.onfido.com/documentation/*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %v\n", t.Token)

	if err := client.DeleteApplicant(ctx, applicant.ID); err != nil {
		panic(err)
	}
}
