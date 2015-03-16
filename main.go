package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"mailchamps/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	timeoutHandler := http.TimeoutHandler(router, time.Second*3, "Timeout!")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	log.Println("Starting server...")

	router.HandleFunc("/subscribe", func(writer http.ResponseWriter, request *http.Request) {

		decoder := json.NewDecoder(request.Body)
		var input api.SubscriptionInput
		err := decoder.Decode(&input)

		if len(input.Email) == 0 || len(input.Username) == 0 {
			var message string
			if len(input.Email) == 0 && len(input.Username) == 0 {
				message = "Incorrect parameters"
			} else if len(input.Email) == 0 {
				message = "Missing email"
			} else if len(input.Username) == 0 {
				message = "Missing username"
			}
			incorrectParametersOutput := api.SubscriptionOutput{
				Email:      input.Email,
				Message:    message,
				StatusCode: 409,
			}

			incorrectParametersOutput.SendSubscriptionResponse(writer)
		} else if err != nil {
			errorOutput := api.SubscriptionOutput{
				Email:      input.Email,
				Message:    "Incorrect parameters",
				StatusCode: 400,
			}

			errorOutput.SendSubscriptionResponse(writer)
		} else {
			api.AddNewsletterSubscriber(writer, request, input)
		}
	}).Methods("POST")

	http.ListenAndServe(":"+port, timeoutHandler)
}
