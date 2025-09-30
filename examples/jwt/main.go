package main

import (
	"context"
	"fmt"

	onfido "github.com/uw-labs/go-onfido"
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
		DOB:       "1990-01-31",
		Location: onfido.Location{ // New mandatory field for v3.4+
			CountryOfResidence: "GBR",
		},
		Address: &onfido.Address{ // Now single address instead of array
			BuildingNumber: "18",
			Street:         "Wind Corner",
			Town:           "Crawley",
			State:          "West Sussex",
			Postcode:       "NW9 5AB",
			Country:        "GBR",
		},
		// For US applicants, consents are mandatory:
		// Consents: []onfido.Consent{
		// 	{
		// 		Name:    string(onfido.ConsentPrivacyNoticesRead),
		// 		Granted: true,
		// 	},
		// },
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
