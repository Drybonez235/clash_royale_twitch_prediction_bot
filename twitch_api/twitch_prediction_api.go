package twitch_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

type Prediction_data_array struct{
	Data []Prediction_meta_data `json:"data"`
	//Pagination any `json:"pagination"`
}

type Prediction_meta_data struct{
	Id string `json:"id"`
	Broadcaster_id string `json:"broadcaster_id"`
	Broadcaster_name string `json:"broadcaster_name"`
	Broadcaster_login string `json:"broadcaster_login"`
	Title string `json:"title"`
	Winning_outcome_id string `json:"winning_outcome_id"`
	Outcomes []Outcome_response
	Prediction_window int `json:"prediction_window"`
	Status string `json:"status"`
	Created_at string `json:"created_at"`
	Locked_at string `json:"locked_at"`
}

type Outcome_response struct{
	Outcome_id string `json:"id"`
	Outcome_title string `json:"title"`
	Users int `json:"users"`
	Channel_points int `json:"channel_points"`
	Top_predictors []Top_predictors `json:"top_predictors"`
	Color string `json:"color"`
}

type Top_predictors struct{
	User_id string `json:"user_id"`
	User_login string `json:"user_login"`
	User_name string `json:"user_name"`
	Channel_points_used int `json:"channel_points_used"`
	Channel_points_won int 	`json:"channel_points_won"`
}

type Prediction_body struct {
	Broadcaster_id    string    `json:"broadcaster_id"`
	Title            string    `json:"title"`
	Outcomes         []Outcome `json:"outcomes"`
	Prediction_window int       `json:"prediction_window"`
}

type Outcome struct {
	Title string `json:"title"`
}

func Start_prediction(twitch_user sqlite.Twitch_user, Env_struct logger.Env_variables) error {
	//Here we are calling the varify function and passing it all the info it needs. You will need a few if statments if it faisls

	// valid, err := Validate_token(twitch_user.Access_token, twitch_user.User_id, twitch_user.Refresh_token)

	// if err!=nil{
	// 	return err
	// }

	// if !valid {
	// 	twitch_user, err  = sqlite.Get_twitch_user("sub", twitch_user.User_id)

	// 	if err != nil{
	// 		return err
	// 	}
	// }
	prediction_json, err := prediction_body(twitch_user.User_id, twitch_user.Display_Name)
	if err!=nil{
		return err
	}	
	client := &http.Client{}
	req, err := http.NewRequest("POST", Env_struct.PREDICTION_URI, bytes.NewBuffer(prediction_json))
	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: Start_prediction CALL: http.NewRequest " + err.Error())
		return err
	}
	bearer := "Bearer " + twitch_user.Access_token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil{
		err = errors.New("FILE: twitch_prediction FUNC: Start_prediction CALL: client.Do " + err.Error())
		return err
	}
	defer resp.Body.Close()
	var Prediction_data_array Prediction_data_array
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: Start_prediction CALL: io.ReadAll " + err.Error())
		return err
	}
	err = json.Unmarshal(body, &Prediction_data_array)
	if err != nil{
		err = errors.New("FILE: twitch_prediction FUNC: Start_prediction CALL: json.Unmarshal " + err.Error())
		return err
	}	
	err = Prediction_response_parser(Prediction_data_array)
	if err != nil{
		return err
	}
	return err
}

func prediction_body(sub string, display_name string) ([]byte, error){
	prediction_text := fmt.Sprintf(`Will %s win the next game?`, display_name)
	body := Prediction_body{
		Broadcaster_id: sub,
		Title: prediction_text,
		Outcomes: []Outcome{
			{Title: "Yes"},
			{Title: "No"},
		},
		Prediction_window: 60,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		err = errors.New("FILE: twitch_prediction FUNC: prediction_body CALL: json.Marshal " + err.Error())
	return jsonData ,err
	}
	return jsonData, nil
}

//THis is a new untested function... Need to make sure it works. It is called to see if we need to make a new prediction OR wait.
func Check_prediction(sub string, bearer string, prediction_id string, Env_struct logger.Env_variables)(string, error){
	fmt.Println("Check predictions fired")
	client := &http.Client{}
	url_quary := url.Values{}
	url_quary.Set("broadcaster_id", sub)
	url_quary.Set("first", "0")
	//If there is an active prediction in my DB, we will search for it here. If not, we don't set the id.
	if prediction_id != ""{
		fmt.Println("Prediction id was not blank and that means we search for a specific one.")
		url_quary.Set("id", prediction_id)
	}
	url_encoded_string := url_quary.Encode()
	check_prediction_url := Env_struct.PREDICTION_URI +"?"+url_encoded_string
	req, err := http.NewRequest("GET", check_prediction_url, nil)// twitch_prediction_uri ,strings.NewReader(url_encoded_string))

	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: Check_prediction CALL: http.NewRequest " + err.Error())
		return "", err
	}

	bearer_string := "Bearer "+ bearer

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)

	resp, err := client.Do(req)
	if err!=nil {
		err = errors.New("FILE: twitch_prediction FUNC: Check_prediction CALL: client.Do " + err.Error())
		return "", err
	}
	if resp.StatusCode != http.StatusOK{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: client.Do " + resp.Status)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: Check_prediction CALL: io.ReadAll " + err.Error())
		return "", err
	}

	var prediction_body Prediction_data_array

	err = json.Unmarshal(body, &prediction_body)
	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: Check_prediction CALL: json.Unmarshal " + err.Error())
		return "", err
	}

	if prediction_body.Data[0].Status == "ACTIVE" || prediction_body.Data[0].Status == "LOCKED"{
		current_prediction, _, err := sqlite.Get_predictions(sub, "ACTIVE")
		if err!=nil{
			return "", err
		}

		if current_prediction == ""{
			return "not_our_prediction", nil
		}

		if prediction_body.Data[0].Id == current_prediction{
				return "our_prediction", nil
			} else {
				return "not_our_prediction", nil
			}
	} else {
		return "no_active_prediction", nil 
	}
}

func End_prediction(prediction_id string, outcome_id string, broadcaster_id string, bearer_token string, Env_struct logger.Env_variables) error{

	if prediction_id == "" || outcome_id == ""{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction BUG: prediction_id or outcome_id was blank")
		return err
	}

	client := &http.Client{}
	requestBody := map[string]interface{}{
		"broadcaster_id":     broadcaster_id,
		"id":                 prediction_id,
		"status":             "RESOLVED",
		"winning_outcome_id": outcome_id,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: json.Marshal" + err.Error())
		return err
	}
	req, err := http.NewRequest("PATCH", Env_struct.PREDICTION_URI, bytes.NewBuffer(jsonData))
	if err!=nil{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: http.NewRequest" + err.Error())
		return err
	}
	bearer_string := "Bearer "+ bearer_token
	req.Header.Set("Authorization",bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err!=nil {
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: client.Do" + err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: client.Do " + resp.Status)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: io.ReadAll" + err.Error())
		return err
	}
	var json_message map[string]interface{} 
	json.Unmarshal(body, &json_message)
	err = sqlite.Delete_prediction_id(broadcaster_id)
	if err !=nil{
		return err
	}
	return nil
}

func Cancel_prediction(prediction_id string, broadcaster_id string, bearer_token string, Env_struct logger.Env_variables)(error){
	if prediction_id == "" {
		err := errors.New("FILE: twitch_prediction FUNC: Cancel_prediction BUG: prediction_id was blank")
		return err
	}

	client := &http.Client{}
	requestBody := map[string]interface{}{
		"broadcaster_id":     broadcaster_id,
		"id":                 prediction_id,
		"status":             "CANCELED",
		"winning_outcome_id": "null",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		err = errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: json.Marshal" + err.Error())
		return err
	}
	req, err := http.NewRequest("PATCH", Env_struct.PREDICTION_URI, bytes.NewBuffer(jsonData))
	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: http.NewRequest" + err.Error())
		return err
	}
	bearer_string := "Bearer "+ bearer_token
	req.Header.Set("Authorization",bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err!=nil {
		err = errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: client.Do" + err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: client.Do " + resp.Status)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		err = errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: io.ReadAll" + err.Error())
		return err
	}
	var json_message map[string]interface{} 
	json.Unmarshal(body, &json_message)
	err = sqlite.Delete_prediction_id(broadcaster_id)
	if err !=nil{
		return err
	}
	return nil

}


func Prediction_response_parser(prediction_data_array Prediction_data_array) error{
	data := prediction_data_array.Data
	prediction := data[0]
	prediction_id := prediction.Id
	streamer_id := prediction.Broadcaster_id
	created_at := prediction.Created_at
	created_at = created_at[0:19]+"Z"
	err := sqlite.Write_new_prediction(streamer_id, prediction_id, created_at)
	if err != nil{
		return err
	}
	var write_outcomes []map[string]interface{} 
	for i := 0; i < len(prediction.Outcomes); i++{
		maps := make(map[string]any)
		outcome := prediction.Outcomes[i]
		maps["prediction_id"] = prediction_id
		maps["outcome_id"] = outcome.Outcome_id
		maps["title"] = outcome.Outcome_title
		if outcome.Outcome_title == "Yes"{
			maps["lose_win"] = 1
		} else {
			maps["lose_win"] = 0	
		}
		write_outcomes = append(write_outcomes, maps)
	}
	err = sqlite.Write_new_prediction_outcomes(write_outcomes)
	if err!=nil{
		return err
	}
	return nil
}