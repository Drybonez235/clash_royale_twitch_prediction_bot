package twitch_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const twitch_prediction_uri = "https://api.twitch.tv/helix/predictions"

//We need a way of getting the dispaly name.

func start_prediction(access_token string, sub string, display_name string) error {

	json, err := prediction_body(sub, display_name)

	if err != nil{
		return err
	}

	//client := &http.Client{}

	req, err := http.NewRequest("GET", twitch_prediction_uri, bytes.NewBuffer(json))

	if err!=nil{
		err = errors.New("there was a problem forming the request")
		return err
	}
	bearer := "bearer " + access_token

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Client-Id", client_id)
	req.Header.Add("Content-Type", "application/json")
	

	//resp, err := client.Do(req)
	// down here we are creating the prediction.
	return err
}

func prediction_body(sub string, display_name string) ([]byte, error){

	jsondata := fmt.Sprintf(`
		{"broadcaster_id": %s,
		"title": fmt.Sprintf("Will %s win the next game?",
		"outcomes":[{
			{"title": "Yes"},
			{"title": "No"},
		}],
		"prediction_window": 60,
	}`,sub, display_name )

	jsonData, err := json.Marshal(jsondata)
	if err != nil {
		err = errors.New("Problem with marshaling the data")
	}

	return jsonData, err
}