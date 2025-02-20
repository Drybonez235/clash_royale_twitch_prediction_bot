package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"fmt"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

//const sub_uri = "https://api.twitch.tv/helix/eventsub/subscriptions"
const sub_uri = "http://localhost:8080/eventsub/subscriptions"

//This is the callback that needs to handle the challenge
const my_website = "http://localhost:3000/subscription_handler"
const app_secret = "dacb3721ea3023f1e955a053d91f24"
const App_id ="b2109dc3a41733acaa7b3fa355df4c" //Test app id

func Create_EventSub(user sqlite.Twitch_user, sub_type string)(error){
	client := http.Client{}

	bearer_string := "Bearer " + app_secret 

	// url_quary := url.Values{}
	// url_quary.Set("Authorization", bearer_string)
	// url_quary.Set("Client-Id", App_id)
	// url_quary.Set("Content-Type", "application/json")
	// url_quary.Encode()

	// fmt.Println(url_quary)

	req_body, err := create_sub_request_body(user, sub_type)


	if err!=nil{
		fmt.Println(err)
		return err
	}
	fmt.Println(string(req_body))

	req, err := http.NewRequest("POST", sub_uri, bytes.NewBuffer(req_body))

	if err!=nil{
		fmt.Println("err")
		return err
	}

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", App_id)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err!=nil || resp.StatusCode != http.StatusOK{
		fmt.Println("Do we get down here? Its the probl")
		fmt.Println(resp.Status)
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.Status)
	defer resp.Body.Close()
	fmt.Println("Do we get down here?")

	return nil
}


func create_sub_request_body(user sqlite.Twitch_user, sub_type string)([]byte, error){
	// The sub types for the CLI and the production API are different
	// 
	// if sub_type == "stream.online" {
	// } else if sub_type == "stream.offline" {
	// } else {
	// 	err := errors.New("invalid sub type requesed hmm")
	// 	return nil, err	
	// }

	//FIX THE TYPES
	err := errors.New("No err")
	
	body := EventSubRequest{
		Type: sub_type,
		Version: "1",
		Condition: ConditionSubRequest{
			UserID: user.User_id,
		},
		Transport : TransportSubRequest{
			Method: "webhook",
			Callback: my_website,
			Secret: secret},
		}

	body_byte, err := json.Marshal(&body)

	if err != nil{
		return nil, err
	}

	return body_byte, nil
}

func respond_challenge(w http.ResponseWriter, body []byte){
	fmt.Println("Respond to challenge fired")
	var data Challenge_struct
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	challenge := data.Challenge
	if challenge ==""{
		http.Error(w, "Challenge not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
}
