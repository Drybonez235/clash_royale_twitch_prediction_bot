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

const redirect_uri = "http://localhost:3000"

type twitch_user_info struct{
	app_request string
	app_received string
	token_exp float64
	token_iat float64
	token_iss string
	sub string
	access_token string
	refresh_token string
	scope string
	token_type string
}

func request_oath_token(code string) (error) {
	twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	url_quary := url.Values{}
		
	url_quary.Set("client_id", client_id)
	url_quary.Set("client_secret", "")
	url_quary.Set("grant_type", "authorization_code")
	url_quary.Set("code", code)
	url_quary.Set("redirect_uri", redirect_uri)
	url_encoded_string := url_quary.Encode()
	
	client := &http.Client{}

	req, err := http.NewRequest("POST", twitch_oath_url, strings.NewReader(url_encoded_string))

	if err != nil {
		err = errors.New("there was something wrong with the POST request")
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		err = errors.New("there was something wrong with the POST Response")
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		err = errors.New("there was something wrong with the POST response body")
		return err
	}

	json_array := make(map[string]any)

	err = json.Unmarshal(body, &json_array)

	if err != nil {
		err = errors.New("there was something wrong with the json response")
		return err	
	}

	var new_user twitch_user_info

	new_user.access_token = json_array["access_token"].(string)
	new_user.refresh_token = json_array["refresh_token"].(string)
	new_user.scope = "channel:manage:predictions"
	new_user.token_type = json_array["token_type"].(string)

	err = Get_claims(new_user.access_token, new_user)

	if err != nil{
		return err
	}

	return err
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
		err = errors.New("there was something wrong with the GET request")
		return "", err	
	}

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		err = errors.New("there was something wrong with the GET response")
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil || resp.Status != "200 OK"{
		err = errors.New("there was something wrong with the GET body response")
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
		err = errors.New("there was something wrong with the GET request")
		return "", err	
	}

	req.Header.Add("OAuth", AOauth_token)

	resp, err := client.Do(req)

	if err != nil{
		err = errors.New("there was something wrong with the GET response")
		return "", err	
	}

	defer resp.Body.Close()

	if err != nil || resp.Status != "200 OK"{
		err = errors.New("there was something wrong with the GET body response")
		return "", err
	} 
	return resp.Status, err
}

// func Get_user_info(streamer_id string, client_secret string){
// 	client_id := "now6dwkymg4vo236ius5d0sn82v9ul"

// 	OAuth_token, err := request_oath_token(client_id, client_secret)
	
// 	if err != nil{
// 		panic(err)
// 	}

// 	valid_token, err := validate_token(OAuth_token)

// 	if err != nil || valid_token != "200 OK"{
// 		err = errors.New("the oauth token is not valid")
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	user_data, err := call_user_endpoint(streamer_id, OAuth_token, client_id)

// 	if err != nil{
// 		panic(err)
// 	}

// 	fmt.Println(user_data)
// }

func Get_claims(oauth_token string, new_user twitch_user_info) (error){
	fmt.Println("Get claims ran")

	twitch_verifiy_user_endpoint := "https://id.twitch.tv/oauth2/userinfo"

	bearer_token := "Bearer " + oauth_token

	fmt.Println(bearer_token)

	client := &http.Client{}
	
	req, err := http.NewRequest("GET", twitch_verifiy_user_endpoint, nil)

	req.Header.Add("Authorization", bearer_token)
	req.Header.Add("Content-Type", "application/json")

	if err != nil{
		fmt.Println(err)
		err = errors.New("there was something wrong with the Get request")
		return err	
	}

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		fmt.Println(err)
		err = errors.New("there was something wrong with the Get response")
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil{
		fmt.Println(err)
		err = errors.New("there was something wrong with the GET body response")
		return err
	}

	json_array := make(map[string]any)

	err = json.Unmarshal(body, &json_array)

	fmt.Println("Did we unmarshel the json?")

	if err != nil {
		err = errors.New("there was something wrong with the json response")
		return err	
	}

	new_user.app_request = json_array["aud"].(string)
	new_user.app_received = json_array["azp"].(string)
	new_user.token_exp = json_array["exp"].(float64)
	new_user.token_iat = json_array["iat"].(float64)
	new_user.token_iss = json_array["iss"].(string)
	new_user.sub = json_array["sub"].(string)

	fmt.Println(new_user)

	return err
}
