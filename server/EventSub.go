package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

func Create_EventSub(user sqlite.Twitch_user, sub_type string, Env_struct logger.Env_variables)(error){
	client := http.Client{}

	bearer_string := "Bearer " + Env_struct.APP_SECRET

	req_body, err := create_sub_request_body(user, sub_type, Env_struct)

	if err!=nil{
		return err
	}
	//Subscription URI info is either "https://api.twitch.tv/helix/eventsub/subscriptions" or "http://localhost:8080/eventsub/subscriptions"
	req, err := http.NewRequest("POST", Env_struct.SUBSCRIPTION_INFO_URI, bytes.NewBuffer(req_body))
	if err!=nil{
		err = errors.New("FILE: EventSub FUNC: Create_EventSub CALL: http.NewRequest " + err.Error())	
		return err
	}

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err!=nil || resp.StatusCode != http.StatusOK{
		err = errors.New("FILE: EventSub FUNC: Create_EventSub CALL: client.Do " + err.Error())	
		return err
	}

	if resp.StatusCode != http.StatusOK{
		err = errors.New("FILE: EventSub FUNC: Create_EventSub CALL: client.DO " + resp.Status)
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
			Secret: Env_struct.APP_SECRET},
		}

	body_byte, err := json.Marshal(&body)

	if err != nil{
		err = errors.New("FILE: EventSub FUNC: create_sub_request_body CALL: json.Marshal " + err.Error())	
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
