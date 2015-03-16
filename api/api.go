//  api.go
//  Mailchamp
//
//  Created by Joe Fabisevich on 3/16/15

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mattbaird/gochimp"
)

////////////////////////////////////////////////////////////////////////////////
// Subscription Input

type SubscriptionInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

////////////////////////////////////////////////////////////////////////////////
// Subscription Output

type SubscriptionOutput struct {
	Email      string `json:"email"`
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}

////////////////////////////////////////////////////////////////////////////////
// SubscriptionOutput

func (output SubscriptionOutput) SendSubscriptionResponse(rw http.ResponseWriter) {
	bytes, err := json.Marshal(output)

	if err != nil {
		log.Fatal("An error occurred marshalling the JSON: ", err)
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(bytes)
}

////////////////////////////////////////////////////////////////////////////////
// API

func AddNewsletterSubscriber(rw http.ResponseWriter, req *http.Request, input SubscriptionInput) {
	apiKey := os.Getenv("MAILCHIMP_API_KEY")

	if len(apiKey) == 0 {
		log.Fatal("Missing Mailchimp API key! Ooh ooh ah ah!")
	}

	api := gochimp.NewChimp(apiKey, true)

	log.Println("Adding subscriber: " + input.Username + "/" + input.Email)

	mergeVars := make(map[string]interface{})
	mergeVars["USERNAME"] = input.Username

	email := gochimp.Email{
		Email: input.Email,
	}

	newsletterId := os.Getenv("NEWSLETTER_ID")

	subscriber := gochimp.ListsSubscribe{
		ApiKey:           apiKey,
		ListId:           newsletterId,
		Email:            email,
		MergeVars:        mergeVars,
		EmailType:        "html",
		DoubleOptIn:      false,
		UpdateExisting:   true,
		ReplaceInterests: true,
		SendWelcome:      false,
	}

	email, err := api.ListsSubscribe(subscriber)

	if err != nil {
		log.Fatal("An error occurred subscribing "+input.Username+" to the list: ", err.Error())
	}

	output := SubscriptionOutput{
		Message:    input.Username + " has been subscribed",
		StatusCode: 200,
		Email:      input.Email,
	}

	output.SendSubscriptionResponse(rw)

	log.Println("Subscriber response sent")
}
