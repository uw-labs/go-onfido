package main

import (
	"fmt"
	"net/http"

	"github.com/tumelohq/go-onfido"
)

func main() {
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
