package twitch_api

import (
	//"bytes"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	//"strings"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

//const twitch_prediction_uri = "https://api.twitch.tv/helix/predictions"
const twitch_prediction_uri = "http://localhost:8080/mock/predictions"

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

func Start_prediction(twitch_user sqlite.Twitch_user) error {
	fmt.Println("Start of prediction functiom")

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
	prediction_json := prediction_body(twitch_user.User_id, twitch_user.Display_Name)
	client := &http.Client{}
	req, err := http.NewRequest("POST", twitch_prediction_uri, bytes.NewBuffer(prediction_json))
	if err!=nil{
		err = errors.New("there was a problem forming the request")
		return err
	}
	bearer := "Bearer " + twitch_user.Access_token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Client-Id", App_id)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println(resp.StatusCode)
		fmt.Println("There was a problem with the response")
		return err
	}
	defer resp.Body.Close()
	var Prediction_data_array Prediction_data_array
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		return err
	}
	err = json.Unmarshal(body, &Prediction_data_array)
	if err != nil{
		return err
	}	
	err = Prediction_response_parser(Prediction_data_array)
	if err != nil{
		return err
	}
	return err
}

func prediction_body(sub string, display_name string) ([]byte){
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
	//body := fmt.Sprintf(`{"broadcaster_id":"%s","title":"%s","outcomes":[{"title":"Yes"},{"title":"No"}],"prediction_window":60}`,sub, prediction_text)
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
	}
	fmt.Println("Made the prediction body")
	return jsonData
}

//THis is a new untested function... Need to make sure it works. It is called to see if we need to make a new prediction OR wait.
func Check_prediction(sub string, bearer string, prediction_id string)(string, error){
	fmt.Println("CHeck Prediction fired")

	fmt.Println("Check Prediction sub: " + sub)
	client := &http.Client{}

	url_quary := url.Values{}
	url_quary.Set("broadcaster_id", sub)
	url_quary.Set("first", "0")

	//If there is an active prediction in my DB, we will search for it here. If not, we don't set the id.
	if prediction_id != ""{
		url_quary.Set("id", prediction_id)
	}

	url_encoded_string := url_quary.Encode()

	check_prediction_url := twitch_prediction_uri +"?"+url_encoded_string

	fmt.Println(check_prediction_url)

	req, err := http.NewRequest("GET", check_prediction_url, nil)// twitch_prediction_uri ,strings.NewReader(url_encoded_string))

	if err!=nil{
		return "", err
	}

	bearer_string := "Bearer "+ bearer

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", App_id)

	resp, err := client.Do(req)

	if err!=nil || resp.StatusCode != http.StatusOK{
		fmt.Println("This is the problem")
		fmt.Println(resp.Status)
		return "", err
	}	

	body, err := io.ReadAll(resp.Body)

	if err!=nil{
		return "", err
	}

	var prediction_body Prediction_data_array

	err = json.Unmarshal(body, &prediction_body)

	if err!=nil{
		return "", err
	}
	fmt.Println(prediction_body.Data[0])

	if prediction_body.Data[0].Status == "ACTIVE" || prediction_body.Data[0].Status == "LOCKED"{
		current_prediction, _, err := sqlite.Get_predictions(sub, "ACTIVE")

		fmt.Println(current_prediction)
		fmt.Println(prediction_body.Data[0].Id)

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

func End_prediction(prediction_id string, outcome_id string, broadcaster_id string, bearer_token string) error{

	if prediction_id == "" || outcome_id == ""{
		err := errors.New("prediction or outcome id was blank")
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
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	req, err := http.NewRequest("PATCH", twitch_prediction_uri, bytes.NewBuffer(jsonData))
	if err!=nil{
		return err
	}
	bearer_string := "Bearer "+ bearer_token
	req.Header.Set("Authorization",bearer_string)
	req.Header.Set("Client-Id", App_id)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err!=nil || resp.StatusCode != http.StatusOK{
		fmt.Println("This is causing the problem")
		fmt.Println(resp.Status)
		fmt.Println(resp.Header)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		return err
	}
	var json_message map[string]interface{} 
	json.Unmarshal(body, &json_message)
	fmt.Println("We closed the prediction")
	err = sqlite.Delete_prediction_id(broadcaster_id)
	if err !=nil{
		return err
	}
	return nil
}

func Cancel_prediction(prediction_id string, broadcaster_id string, bearer_token string)(error){
	fmt.Println("Cancel prediction fired")
	if prediction_id == "" {
		err := errors.New("prediction id was blank")
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
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	req, err := http.NewRequest("PATCH", twitch_prediction_uri, bytes.NewBuffer(jsonData))
	if err!=nil{
		return err
	}
	bearer_string := "Bearer "+ bearer_token
	req.Header.Set("Authorization",bearer_string)
	req.Header.Set("Client-Id", App_id)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err!=nil || resp.StatusCode != http.StatusOK{
		fmt.Println("This is causing the problem")
		fmt.Println(resp.Status)
		fmt.Println(resp.Header)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		return err
	}
	var json_message map[string]interface{} 
	json.Unmarshal(body, &json_message)
	fmt.Println("We canceled the prediction")
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