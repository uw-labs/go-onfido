package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/uw-labs/go-onfido"
)

func main() {
	wh, err := onfido.NewWebhookFromEnv()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/webhook/onfido", func(w http.ResponseWriter, req *http.Request) {
		whReq, err := wh.ParseFromRequest(req)
		if err != nil {
			if errors.Is(err, onfido.ErrInvalidWebhookSignature) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid signature")) //nolint:errcheck

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error occurred")) //nolint:errcheck

			return
		}

		fmt.Fprintf(w, "Webhook: %+v\n", whReq) //nolint:errcheck
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
