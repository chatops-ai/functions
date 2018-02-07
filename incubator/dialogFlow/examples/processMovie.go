import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
)

func ProcessMessage(event Messaging) {
	var userQuery = event.Message.Text
	var dialogFlowResponse = GetResponse(userQuery)
	client := &http.Client{}
	var response Response

	if !reflect.DeepEqual(dialogFlowResponse.Metadata, apiai.Metadata{}) && dialogFlowResponse.Metadata.IntentName == "shows" {
		var showType = dialogFlowResponse.Parameters["show-type"]
		db := NewMovieDB()
		var shows []Show

		if showType == "movie" {
			shows = db.GetNowPlayingMovies()
		} else {
			shows = db.GetAiringTodayShows()
		}

		response = Response{
			Recipient: User{
				ID: event.Sender.ID,
			},
			Message: Message{
				Attachment: &Attachment{
					Type: "template",
					Payload: Payload{
						TemplateType: "generic",
						Elements:     BuildCarousel(shows[:10]),
					},
				},
			},
		}
	} else {
		response = Response{
			Recipient: User{
				ID: event.Sender.ID,
			},
			Message: Message{
				Text: dialogFlowResponse.Fulfillment.Speech,
			},
		}
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)

	url := fmt.Sprintf(FACEBOOK_API, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
