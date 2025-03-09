package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger" 
)

func Create_EventSub(user sqlite.Twitch_user, sub_type string, Env_struct logger.Env_variables)(error){
	client := http.Client{}

	bearer_string := "Bearer " + Env_struct.APP_SECRET

	// url_quary := url.Values{}
	// url_quary.Set("Authorization", bearer_string)
	// url_quary.Set("Client-Id", App_id)
	// url_quary.Set("Content-Type", "application/json")
	// url_quary.Encode()

	// fmt.Println(url_quary)

	req_body, err := create_sub_request_body(user, sub_type, Env_struct)

	if err!=nil{
		fmt.Println(err)
		return err
	}
	fmt.Println(string(req_body))

	//Subscription URI info is either "https://api.twitch.tv/helix/eventsub/subscriptions" or "http://localhost:8080/eventsub/subscriptions"

	req, err := http.NewRequest("POST", Env_struct.SUBSCRIPTION_INFO_URI, bytes.NewBuffer(req_body))

	if err!=nil{
		fmt.Println("err")
		return err
	}

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err!=nil || resp.StatusCode != http.StatusOK{
		return err
	}
	defer resp.Body.Close()

	return nil
}


func create_sub_request_body(user sqlite.Twitch_user, sub_type string, Env_struct logger.Env_variables)([]byte, error){
	// The sub types for the CLI and the production API are different
	// 
	// if sub_type == "stream.online" {
	// } else if sub_type == "stream.offline" {
	// } else {
	// 	err := errors.New("invalid sub type requesed hmm")
	// 	return nil, err	
	// }

	callback_string :=  Env_struct.ROYALE_BETS_URL+"/subscription_handler"
	
	body := EventSubRequest{
		Type: sub_type,
		Version: "1",
		Condition: ConditionSubRequest{
			UserID: user.User_id,
		},
		Transport : TransportSubRequest{
			Method: "webhook",
			Callback: callback_string,
			Secret: secret},
		}

	body_byte, err := json.Marshal(&body)

	if err != nil{
		return nil, err
	}

	return body_byte, nil
}

func respond_challenge(w http.ResponseWriter, body []byte){
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
