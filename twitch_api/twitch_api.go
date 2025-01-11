package twitch_api

import (
	"strings"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func request_oath_token(client_id string, client_secret string) (string, error) {
	twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	url_quary := url.Values{}
		
	url_quary.Set("client_id", client_id)
	url_quary.Set("client_secret", client_secret)
	url_quary.Set("grant_type", "client_credentials")
	url_encoded_string := url_quary.Encode()
	
	client := &http.Client{}

	req, err := http.NewRequest("POST", twitch_oath_url, strings.NewReader(url_encoded_string))

	if err != nil {
		err = errors.New("There was something wrong with the POST request")
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		err = errors.New("There was something wrong with the POST Response")
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		err = errors.New("There was something wrong with the POST response body")
		return "", err
	}

	json_array := make(map[string]any)

	err = json.Unmarshal(body, &json_array)

	if err != nil {
		err = errors.New("There was something wrong with the json response")
		return "", err	
	}
	access_token_string := json_array["access_token"].(string)

	return access_token_string, err
}

func call_user_endpoint(streamer_id string, access_token string, client_id string) (string, error) {
	twitch_get_user_endpoint := "https://api.twitch.tv/helix/users?"
	bearer_token := "Bearer " + access_token

	url_quary := url.Values{}
	url_quary.Set("login", streamer_id)
	
	url_encoded_string := url_quary.Encode()

	get_call := twitch_get_user_endpoint + url_encoded_string
	
	client := &http.Client{}
	
	req, err := http.NewRequest("GET", get_call, nil)
	req.Header.Add("Authorization", bearer_token)
	req.Header.Add("Client-Id", client_id)

	if err != nil{
		err = errors.New("There was something wrong with the GET request")
		return "", err	
	}

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		err = errors.New("There was something wrong with the GET response")
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil || resp.Status != "200 OK"{
		err = errors.New("There was something wrong with the GET body response")
		return "", err
	}

	json_response := string(body)

	return json_response, err
}

func validate_token(AOauth_token string) (string, error){
	twitch_validation_endpoint := "https://id.twitch.tv/oauth2/validate"

	client := &http.Client{}

	req, err := http.NewRequest("GET", twitch_validation_endpoint, nil)

	if err != nil{
		err = errors.New("There was something wrong with the GET request")
		return "", err	
	}

	req.Header.Add("OAuth", AOauth_token)

	resp, err := client.Do(req)

	if err != nil{
		err = errors.New("There was something wrong with the GET response")
		return "", err	
	}

	defer resp.Body.Close()

	//body, err := io.ReadAll(resp.Body)

	if err != nil || resp.Status != "200 OK"{
		err = errors.New("There was something wrong with the GET body response")
		return "", err
	} 
	return resp.Status, err
}

func Get_user_info(streamer_id string, client_secret string){
	client_id := "now6dwkymg4vo236ius5d0sn82v9ul"

	OAuth_token, err := request_oath_token(client_id, client_secret)
	
	if err != nil{
		panic(err)
	}

	valid_token, err := validate_token(OAuth_token)

	if err != nil || valid_token != "200 OK"{
		err = errors.New("The OAuth token is not valid.")
		fmt.Println(err)
		panic(err)
	}

	user_data, err := call_user_endpoint(streamer_id, OAuth_token, client_id)

	if err != nil{
		panic(err)
	}

	fmt.Println(user_data)
} 
