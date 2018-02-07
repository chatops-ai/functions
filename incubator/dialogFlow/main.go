package main

import (
	"log"
	"os"

	. "github.com/mlabouardy/dialogflow-go-client"
	. "github.com/mlabouardy/dialogflow-go-client/models"
)

// GetResponse queries DialogFlow with the users payload
// in order to parse a reponse via the DialogFlow agent
func ParseInput(input string) Result {
	err, client := NewDialogFlowClient(Options{
		AccessToken: os.Getenv("DIALOG_FLOW_TOKEN"),
	})
	if err != nil {
		log.Fatal(err)
	}

	query := Query{
		Query: input,
	}
	resp, err := client.QueryFindRequest(query)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Result
}
