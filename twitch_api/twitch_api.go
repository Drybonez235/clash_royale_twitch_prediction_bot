package twitch_api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

func request_oath_token(client_id string, secret string) (error) {
	twitch_oath_url := "https://id.twitch.tv/oauth2/token"
	
	client := &http.Client{}
	
	req, err := http.NewRequest("Post", twitch_oath_url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		fmt.Print(err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		fmt.Print(err)
		return err
	}

	access_token_json, err := data_handler_array(body)

	if err != nil {
		fmt.Print(err)
		return err	
	}
}

func data_handler_array(data []byte) ([3]string, error){
	var json_array []map[string]any
	var access_token_array [3]string

	err := json.Unmarshal(data, &json_array)
	
	if(err != nil){
		err = errors.New("something went wrong")
	}
	zero := json_array[0] 
	
	access_token_array[0] =  zero["access_token"].(string)
	access_token_array[1] =  zero["expires_in"].(string)	
	access_token_array[2] =  zero["token_type"].(string)
 
	return access_token_array, err
}

func Get_user

func Get_twitch() string {
	name := "testing_my_first_app"
	oath_redirect := "http://localhost:3000"
	clien_id := "now6dwkymg4vo236ius5d0sn82v9ul"

	return "Hello from twitch"
}