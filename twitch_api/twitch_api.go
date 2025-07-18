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
	"strings"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/ncruces/go-sqlite3"
)

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

type EventSubRequest struct {
	Type      string     `json:"type"`
	Version   string     `json:"version"`
	Condition ConditionSubRequest  `json:"condition"`
	Transport TransportSubRequest  `json:"transport"`
}

type ConditionSubRequest struct {
	UserID string `json:"user_id"`
}

type TransportSubRequest struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

//This function takes the code in the response from the twitch streamer granting access and obtains an OAuth Token.
//It sends a POST request with the data in the URL encoded body (NOT to be confused with query parameters). form encoded url encoded parameters.
//The twitch API infrmation is found at https://dev.twitch.tv/docs/authentication/getting-tokens-oidc/#oidc-authorization-code-grant-flow
func Request_user_oath_token(code string, player_tag string, Env_struct logger.Env_variables, db *sqlite3.Conn) (error) {
	fmt.Println("Request_user_oath_token ran")
	//twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	url_quary := url.Values{}
	url_quary.Set("client_id", Env_struct.APP_ID)
	url_quary.Set("client_secret", Env_struct.APP_SECRET)
	url_quary.Set("grant_type", "authorization_code")
	url_quary.Set("code", code)
	url_quary.Set("redirect_uri", Env_struct.ROYALE_BETS_URL+"/redirect")
	url_encoded_string := url_quary.Encode()
	client := &http.Client{}
	req, err := http.NewRequest("POST", Env_struct.OAUTH_REFRESH_TOKEN_URI +"?", strings.NewReader(url_encoded_string))
	if err != nil {
		err = errors.New("FILE: twitch_api FUNC: Generate_state_nonce CALL: http.NewRequest" + err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil{
		fmt.Print(err)
		err = errors.New("FILE: twitch_api FUNC: Generate_state_nonce CALL: client.DO " + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: Generate_state_nonce CALL: io.ReadAll" + err.Error())
		return err
	}
	var new_user access_token_response_json
	err = json.Unmarshal(body, &new_user)
	if err != nil {
		err = errors.New("FILE: twitch_api FUNC: Generate_state_nonce CALL: json.Unmarshal " + err.Error())
		return err	
	}
	err = Get_claims(new_user.Access_token, new_user, player_tag, Env_struct, db)
	if err != nil{
		return err
	}
	return err
}

//Gets the claims information assoicated with an OAth token. This is where we get the sub ID.
//Sends a GET request to https://id.twitch.tv/oauth2/userinfo with the user access token as a Bearer header.
//Information about the Twitch API can be found at: https://dev.twitch.tv/docs/authentication/getting-tokens-oidc/#getting-claims-information-from-an-access-token
func Get_claims(oauth_token string, new_user access_token_response_json, player_tag string, Env_struct logger.Env_variables, db *sqlite3.Conn) (error){
	//twitch_verifiy_user_endpoint := "https://id.twitch.tv/oauth2/userinfo"
	fmt.Println("Get Claims Ran")
	bearer_token := "Bearer " + oauth_token
	client := &http.Client{}
	req, err := http.NewRequest("GET", Env_struct.OAUTH_CLAIMS_INFO_URI, nil)
	req.Header.Add("Authorization", bearer_token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: Get_claims CALL: http.NewRequest " + err.Error())
		return err	
	}
	resp, err := client.Do(req)
	if err != nil || resp.Status != "200 OK"{
		err = errors.New("FILE: twitch_api FUNC: Get_claims CALL: client.DO " + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: Get_claims CALL: io.ReadAll " + err.Error())
		return err
	}
	var claims_json claims_json
	err = json.Unmarshal(body, &claims_json)
	if err != nil {
		err = errors.New("FILE: twitch_api FUNC: Get_claims CALL: io.Unmarshal " + err.Error())
		return err	
	}
	display_name, err := Get_display_name(claims_json.Sub, Env_struct)
	if err!=nil{
		return err
	}
	unpacked_scope, err := Scope_unpacker(new_user.Scope)
	if err != nil{
		return err
	}

	err = sqlite.Write_twitch_info(db ,claims_json.Sub, display_name, new_user.Access_token, new_user.Refresh_token, unpacked_scope, new_user.Token_type,
		claims_json.Aud, claims_json.Aud, claims_json.Exp, claims_json.Iat, claims_json.Iss, 0, player_tag)

	if err!=nil{
		return err
	}
	//fmt.Println(claims_json.Sub)
	
	//So EVENTsub sends a request to twitch to sub to a streamer. When you request this, twitch sends back a response BUT it is required that the response has to be sent to a website using HTTPS and on port 443. More info here: https://dev.twitch.tv/docs/api/reference/#create-eventsub-subscription
	err = Create_EventSub(claims_json.Sub, "stream.online", Env_struct)
	if err!=nil{return err}

	err = Create_EventSub(claims_json.Sub, "stream.offline", Env_struct)
	if err!=nil{return err}

	return nil
}

//This gets the display name of the twitch streamer using the sub id.
//Sends a GET request to https://api.twitch.tv/helix/users and query parameters the data sub id.
//Information about the twitch API can be found at: https://dev.twitch.tv/docs/api/reference/#get-users
func Get_display_name(streamer_id string, Env_struct logger.Env_variables) (string, error) {	
	app_oath_token, err := Request_app_oath_token(Env_struct)
	if err != nil{
		return "", err
	}
	//twitch_get_user_endpoint := "https://api.twitch.tv/helix/users?"
	bearer_token := "Bearer " + app_oath_token
	url_quary := url.Values{}
	url_quary.Set("id", streamer_id)
	url_encoded_string := url_quary.Encode()
	get_call :=  Env_struct.USER_INFO_URI + "?" + url_encoded_string
	client := &http.Client{}
	req, err := http.NewRequest("GET", get_call, nil)
	req.Header.Add("Authorization", bearer_token)
	req.Header.Add("Client-Id", Env_struct.APP_ID)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: Get_display_name CALL: http.NewRequest " + err.Error())
		return "", err	
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		err = errors.New("FILE: twitch_api FUNC: Get_display_name CALL: client.Do " + err.Error())
		return "", err
	}
	defer resp.Body.Close()
	var user_data_array_json user_data_array_json
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: Get_display_name CALL: io.ReadAll " + err.Error())
		return "", err
	}
	err = json.Unmarshal(body, &user_data_array_json)
	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Get_display_name CALL: json.Unmarshal " + err.Error())
		return "", err
	}
	data := user_data_array_json.Data
	if len(data) == 0 {
		err = errors.New("FILE: twitch_api FUNC: Get_display_name BUG: Data was blank")
		return "", err 
	}
	return data[0].DisplayName, nil
}

//This gets an OAuth token for the app.
//Sends a POST request with a URL-Encoded BODY to https://id.twitch.tv/oauth2/token
//Information about the twitch API can be found at: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#client-credentials-grant-flow
func Request_app_oath_token(Env_struct logger.Env_variables) (string, error) {
	url_quary := url.Values{}
	url_quary.Set("client_id", Env_struct.APP_ID)
	url_quary.Set("client_secret", Env_struct.APP_SECRET)
	url_quary.Set("grant_type", "client_credentials")
	url_encoded_string := url_quary.Encode()
	client := &http.Client{}
	req, err := http.NewRequest("POST", Env_struct.OAUTH_REFRESH_TOKEN_URI, strings.NewReader(url_encoded_string))
	if err != nil {
		err = errors.New("FILE: twitch_api FUNC: request_app_oath_token CALL: http.NewRequest " + err.Error())
		return "" ,err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: request_app_oath_token CALL: client.Do " + err.Error())
		return "", err
	}

	if resp.StatusCode != http.StatusOK{
		err := errors.New("FILE: twitch_prediction FUNC: End_prediction CALL: client.Do " + resp.Status)
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: request_app_oath_token CALL: io.ReadAll " + err.Error())
		return "", err
	}
	var app_oauth_token_json app_oauth_token_json 
	err = json.Unmarshal(body, &app_oauth_token_json)
	if err != nil {
		err = errors.New("FILE: twitch_api FUNC: request_app_oath_token CALL: json.Unmarhsal" + err.Error())
		return "", err	
	}
	return app_oauth_token_json.Access_token, nil
}

//This get a new OAuth token using a refresh token for a twitch user who has granted access.
//Sends a POST request with a URL-Encoded BODY to https://id.twitch.tv/oauth2/token
//Information about the twitch API can be found at: https://dev.twitch.tv/docs/authentication/refresh-tokens/#how-to-use-a-refresh-token
func Refresh_token(refresh_token string, user_id string, Env_struct logger.Env_variables, db *sqlite3.Conn) (bool, error){
	client := &http.Client{}
	url_quary := url.Values{}
	url_quary.Set("client_id", Env_struct.APP_ID)
	url_quary.Set("client_secret", Env_struct.APP_SECRET)
	url_quary.Set("grant_type", "refresh_token")
	url_quary.Set("refresh_token", refresh_token)
	url_encoded_string := url_quary.Encode()
	req, err := http.NewRequest("POST", Env_struct.OAUTH_REFRESH_TOKEN_URI, strings.NewReader(url_encoded_string)) 
	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Refresh_token http.NewRequest" + err.Error())
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Refresh_token CALL: client.D" + err.Error())
		return false, err
	}
	json_response, err := io.ReadAll(resp.Body)
	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Refresh_token CALL: io.ReadAll" + err.Error())
		return false, err
	}
	var refresh_token_response Refresh_token_response
	err = json.Unmarshal(json_response, &refresh_token_response)
	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Refresh_token CALL: json.Unmarshal" + err.Error())
		return false, err
	}
	err = sqlite.Update_tokens(db, refresh_token_response.Access_token, refresh_token_response.Refresh_token, user_id)
	if err!=nil{
		return false, err
	}
	return true, nil
}

//Subscribes to events from twitch streamers.
//Sends a POST request using JSON in the body to https://api.twitch.tv/helix/eventsub/subscriptions. 
//Information about the twitch API can be found at: https://dev.twitch.tv/docs/eventsub/manage-subscriptions/#subscribing-to-events and https://dev.twitch.tv/docs/api/reference/#create-eventsub-subscription.
func Create_EventSub(sub_id string, sub_type string, Env_struct logger.Env_variables)(error){
	fmt.Println("Create event sub fired")
	client := http.Client{}
	app_token, err := Request_app_oath_token(Env_struct)
	if err!=nil{return err}
	fmt.Println("Fired before bearer string")
	bearer_string := "Bearer " + app_token //Changed from using the app secret to using an App OAuth token.
	fmt.Println("Fired after bearer string")
	req_body, err := create_sub_request_body(sub_id, sub_type, Env_struct)
	fmt.Println("Fired after making the request body")
	fmt.Println(string(req_body))
	if err!=nil{
		return err
	}
	//Subscription URI info is either "https://api.twitch.tv/helix/eventsub/subscriptions" or "http://localhost:8080/eventsub/subscriptions"
	req, err := http.NewRequest("POST", Env_struct.SUBSCRIPTION_INFO_URI, bytes.NewBuffer(req_body))
	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Create_EventSub CALL: http.NewRequest " + err.Error())	
		return err
	}

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("Fired before making the client do")
	resp, err := client.Do(req)

	if err!=nil {
		err = errors.New("FILE: twitch_api FUNC: Create_EventSub CALL: client.Do " + err.Error())	
		return err
	}

	if resp == nil {
		return errors.New("FILE: EventSub FUNC: Create_EventSub CALL: client.Do returned nil response")
	}

	if resp.StatusCode != http.StatusOK{
		err = errors.New("FILE: EventSub FUNC: Create_EventSub CALL: client.DO " + resp.Status)
		return err
	}
	defer resp.Body.Close()

	return nil
}

//Creates the body used in the POST request for subscribing to events.
//Information about the twitch API can be found at: https://dev.twitch.tv/docs/api/reference/#create-eventsub-subscription
func create_sub_request_body(user_id string, sub_type string, Env_struct logger.Env_variables)([]byte, error){
	callback_string :=  Env_struct.ROYALE_BETS_URL+"/subscription_handler"
	
	body := EventSubRequest{
		Type: sub_type,
		Version: "1",
		Condition: ConditionSubRequest{
			UserID: user_id,
		},
		Transport : TransportSubRequest{
			Method: "webhook",
			Callback: callback_string,
			Secret: Env_struct.ENCRYPTION_SECRET},
		}
	body_byte, err := json.Marshal(&body)
	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: create_sub_request_body CALL: json.Marshal " + err.Error())
		return nil, err
	}
	return body_byte, nil
}