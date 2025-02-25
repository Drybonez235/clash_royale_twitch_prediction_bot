package twitch_api

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	
	"net/http"
	"net/url"
	"strings"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

const App_id = "b2109dc3a41733acaa7b3fa355df4c"
const Secret = "dacb3721ea3023f1e955a053d91f24"//Test secret

const redirect_uri = "http://localhost:3000/redirect"

type access_token_response_json struct {
	Access_token string `json:"access_token"`
	Expires_in int `json:"expires_in"`
	//Id_token string `json:"id_token"`
	Refresh_token string `json:"refresh_token"`
	Scope []string `json:"scope"`
	Token_type string `json:"token_type"`
}

type claims_json struct {
	Aud string `json:"aud"`
	Exp int `json:"exp"`
	Iat int `json:"iat"`
	Iss string `json:"iss"` 
	Sub string `json:"sub"`
}

type app_oauth_token_json struct {
	Access_token string `json:"access_token"`
	Expires_in int `json:"expires_in"`
	Token_type string `json:"token_type"`
}

type user_data_array_json struct {
	Data []user_data_json `json:"data"`
}

type user_data_json struct {
	BroadcasterType  string `json:"broadcaster_type"`
	CreatedAt        string `json:"created_at"`
	Description      string `json:"description"`
	DisplayName      string `json:"display_name"`
	ID               string `json:"id"`
	Login            string `json:"login"`
	OfflineImageURL  string `json:"offline_image_url"`
	ProfileImageURL  string `json:"profile_image_url"`
	Type             string `json:"type"`
	ViewCount        int    `json:"view_count"`
}

type Refresh_token_response struct{
	Access_token string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Scope []string `json:"scope"`
	Token_type string `json:"token_type"`

}

func Request_user_oath_token(code string) (error) {
	fmt.Println("Request_user_oath_token ran")
	//twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	twitch_oath_url := "http://localhost:8080/auth/authorize"
	url_quary := url.Values{}
	url_quary.Set("client_id", App_id)
	url_quary.Set("client_secret", Secret)
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
	if err != nil{
		fmt.Print(err)
		err = errors.New("there was something wrong with the POST Response")
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("there was something wrong with the POST response body")
		return err
	}
	var new_user access_token_response_json
	err = json.Unmarshal(body, &new_user)
	if err != nil {
		err = errors.New("there was something wrong with the json response")
		return err	
	}
	err = Get_claims(new_user.Access_token, new_user)
	if err != nil{
		return err
	}
	return err
}

func Get_claims(oauth_token string, new_user access_token_response_json) (error){
	fmt.Println("Get claims ran")
	//twitch_verifiy_user_endpoint := "https://id.twitch.tv/oauth2/userinfo"
	twitch_verifiy_user_endpoint := "http://localhost:8080/mock/userinfo"
	bearer_token := "Bearer " + oauth_token
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
	var claims_json claims_json
	err = json.Unmarshal(body, &claims_json)
	if err != nil {
		err = errors.New("there was something wrong with the json response")
		return err	
	}
	// display_name, err := Get_display_name(claims_json.Sub)
	display_name := "Mock API Client"
	// if err!=nil{
	// 	return err
	// }
	unpacked_scope, err := Scope_unpacker(new_user.Scope)
	if err != nil{
		return err
	}
//Scope is an array of strings. We need it to be a regular array.
	err = sqlite.Write_twitch_info(claims_json.Sub, display_name, new_user.Access_token, new_user.Refresh_token, unpacked_scope, new_user.Token_type,
		claims_json.Aud, claims_json.Aud, claims_json.Exp, claims_json.Iat, claims_json.Iss, 0, "")
	return err
}

func Get_display_name(streamer_id string) (string, error) {	
	app_oath_token, err := request_app_oath_token()
	if err != nil{
		return "", err
	}
	//twitch_get_user_endpoint := "https://api.twitch.tv/helix/users?"
	twitch_get_user_endpoint := "'http://localhost:8080/mock/users"
	bearer_token := "Bearer " + app_oath_token
	url_quary := url.Values{}
	url_quary.Set("id", streamer_id)
	url_encoded_string := url_quary.Encode()
	get_call := twitch_get_user_endpoint + url_encoded_string
	client := &http.Client{}
	req, err := http.NewRequest("GET", get_call, nil)
	req.Header.Add("Authorization", bearer_token)
	req.Header.Add("Client-Id", App_id)
	if err != nil{
		err = errors.New("there was something wrong with the GET request")
		return "", err	
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		err = errors.New("there was something wrong with the GET response")
		return "", err
	}
	defer resp.Body.Close()
	var user_data_array_json user_data_array_json
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("there was something wrong with the GET body response")
		return "", err
	}
	err = json.Unmarshal(body, &user_data_array_json)
	if err!=nil{
		fmt.Printf("Failed to unmarshal JSON: %v\n", string(body))
		return "", err
	}
	data := user_data_array_json.Data
	if len(data) == 0 {
		return "", errors.New("no data found for the given streamer ID")
	}
	return data[0].DisplayName, nil
}

func request_app_oath_token() (string, error) {
	fmt.Println("Request app _oath_token ran")
	//twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	twitch_oath_url := "http://localhost:8080/units/clients?"
	url_quary := url.Values{}
	url_quary.Set("client_id", App_id)
	url_quary.Set("client_secret", Secret)
	url_quary.Set("grant_type", "client_credentials")
	url_encoded_string := url_quary.Encode()
	client := &http.Client{}
	req, err := http.NewRequest("POST", twitch_oath_url, strings.NewReader(url_encoded_string))
	if err != nil {
		err = errors.New("there was something wrong with the POST request")
		return "" ,err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil || resp.Status != "200 OK"{
		err = errors.New("there was something wrong with the POST Response")
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("there was something wrong with the POST response body")
		return "", err
	}
	var app_oauth_token_json app_oauth_token_json 
	err = json.Unmarshal(body, &app_oauth_token_json)
	if err != nil {
		err = errors.New("there was something wrong with the json response")
		return "", err	
	}
	return app_oauth_token_json.Access_token, err
}

func Validate_token(AOauth_token string, sub string, refresh_token string) (bool, error){
	fmt.Println("Validate token ran")
	twitch_validation_endpoint := "https://id.twitch.tv/oauth2/validate"
	client := &http.Client{}
	req, err := http.NewRequest("GET", twitch_validation_endpoint, nil)
	if err != nil{
		err = errors.New("there was something wrong with the GET request")
		return false, err	
	}
	req.Header.Set("Authorization", "OAuth " + AOauth_token)
	resp, err := client.Do(req)
	if err != nil{
		err = errors.New("there was something wrong with the GET response")
		return false, err	
	}
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		refreshed, err := Refresh_token(refresh_token, sub)
		if err !=nil{
			return false, err
		}
		if !refreshed {
			panic(err)
		}
		return false, nil
	}
	return true, err
}

func Refresh_token(refresh_token string, user_id string) (bool, error){
	refresh_token_url := "https://id.twitch.tv/oauth2/token"
	client := &http.Client{}
	url_quary := url.Values{}
	url_quary.Set("client_id", App_id)
	url_quary.Set("client_secret", Secret)
	url_quary.Set("grant_type", "refresh_token")
	url_quary.Set("refresh_token", refresh_token)
	url_encoded_string := url_quary.Encode()
	req, err := http.NewRequest("POST", refresh_token_url, strings.NewReader(url_encoded_string)) 
	if err!=nil{
		fmt.Println(err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err!=nil{
		return false, err
	}
	json_response, err := io.ReadAll(resp.Body)
	if err!=nil{
		return false, err
	}
	var refresh_token_response Refresh_token_response
	err = json.Unmarshal(json_response, &refresh_token_response)
	if err!=nil{
		return false, err
	}
	err = sqlite.Update_tokens(refresh_token_response.Access_token, refresh_token_response.Refresh_token, user_id)
	if err!=nil{
		return false, err
	}
	return true, nil
}

func Test_request_user_oath_token(user_id string) (error) {
	fmt.Println("TEST Request_user_oath_token ran")
	//twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	twitch_oath_url := "http://localhost:8080/auth/authorize"
	url_quary := url.Values{}
	url_quary.Set("client_id", App_id)
	url_quary.Set("client_secret", Secret)
	url_quary.Set("grant_type", "user_token")
	url_quary.Set("user_id", user_id)
	url_quary.Set("scope", "channel:manage:predictions") //openid") openid
	client := &http.Client{}
	full_url := fmt.Sprintf("%s?%s", twitch_oath_url, url_quary.Encode())
	req, err := http.NewRequest("POST", full_url, nil)
	if err != nil {
		err = errors.New("there was something wrong with the POST request")
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil{
		fmt.Print(err)
		err = errors.New("there was something wrong with the POST Response")
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("there was something wrong with the POST response body")
		return err
	}
	var new_user access_token_response_json
	err = json.Unmarshal(body, &new_user)
	if err != nil {
		err = errors.New("there was something wrong with the json response")
		return err	
	}
	fmt.Println(new_user.Access_token)
	//err = Get_claims(new_user.Access_token, new_user)
	if err != nil{
		return err
	}
	return err
}

func Check_stream_status(sub string) (bool, error){
	//Code that checks to see if the streamer is still streaming.
	return true, nil
}