package twitch_api

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

const twitch_prediction_uri = "https://api.twitch.tv/helix/predictions"

type Prediction_data_array struct{
	Data []string `json:"Data"`
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

	valid, err := Validate_token(twitch_user.Access_token, twitch_user.User_id)

	if !valid{
		//Reresh token refreshes the token and the updates the user
		fmt.Println(twitch_user.Access_token)
		statusOK, err := Refresh_token(twitch_user.Refresh_token, twitch_user.User_id)
		if !statusOK {
			return err
		}

		twitch_user, err = sqlite.Get_twitch_user("sub", twitch_user.User_id)
		fmt.Println(twitch_user.Access_token)
		if err!=nil{
			return err
		}
	}

	prediction_json := prediction_body(twitch_user.User_id, twitch_user.Display_Name)

	if err != nil{
		fmt.Println("We have the json, now what.")
		return err
	}

	client := &http.Client{}

	fmt.Println(prediction_json)

	req, err := http.NewRequest("POST", twitch_prediction_uri, strings.NewReader(prediction_json))

	if err!=nil{
		err = errors.New("there was a problem forming the request")
		return err
	}

	bearer := "Bearer " + twitch_user.Access_token

	req.Header.Set("Authorization", bearer)
	req.Header.Set("Client-Id", App_id)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Right before sending the req")
	fmt.Println(bearer)
	

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK{
		fmt.Println(resp.StatusCode)
		fmt.Print(io.ReadAll(resp.Body))
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

	fmt.Println(Prediction_data_array)

	return err
}

func prediction_body(sub string, display_name string) (string){


	prediction_text := fmt.Sprintf(`Will %s win the next game?`, display_name)

	// body := Prediction_body{
	// 	Broadcaster_id: sub,
	// 	Title: prediction_text,
	// 	Outcomes: []Outcome{
	// 		{Title: "Yes"},
	// 		{Title: "No"},
	// 	},
	// 	Prediction_window: 60,
	// }
	body := fmt.Sprintf(`{"broadcaster_id":"%s","title":"%s","outcomes":[{"title":"Yes"},{"title":"No"}],"prediction_window":60}`,sub, prediction_text)

	//jsonData, err := json.Marshal(body)

	// if err != nil {
	// 	fmt.Println("HERE IS THE PROBLEM")
	// 	err = errors.New("problem with marshaling the data")
	// }
	fmt.Println("Made the prediction body")
	return body
}