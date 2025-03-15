package twitch_api

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
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


func Request_user_oath_token(code string, player_tag string, Env_struct logger.Env_variables) (error) {
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
	req, err := http.NewRequest("POST", Env_struct.OAUTH_REFRESH_TOKEN_URI, strings.NewReader(url_encoded_string))
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
	err = Get_claims(new_user.Access_token, new_user, player_tag, Env_struct)
	if err != nil{
		return err
	}
	return err
}

//Get claims takes an oath token from a twitch streamer oath response and then gets the claims associated with that token. It also relays the player tag associated with state used to verify the oath request.
func Get_claims(oauth_token string, new_user access_token_response_json, player_tag string, Env_struct logger.Env_variables) (error){
	//twitch_verifiy_user_endpoint := "https://id.twitch.tv/oauth2/userinfo"
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
//Scope is an array of strings. We need it to be a regular array.
	err = sqlite.Write_twitch_info(claims_json.Sub, display_name, new_user.Access_token, new_user.Refresh_token, unpacked_scope, new_user.Token_type,
		claims_json.Aud, claims_json.Aud, claims_json.Exp, claims_json.Iat, claims_json.Iss, 0, player_tag)

	if err!=nil{
		return err
	}

	err = Create_EventSub(claims_json.Sub, "stream.online", Env_struct)
	if err!=nil{return err}

	err = Create_EventSub(claims_json.Sub, "stream.offline", Env_struct)
	if err!=nil{return err}

	return nil
}

func Get_display_name(streamer_id string, Env_struct logger.Env_variables) (string, error) {	
	app_oath_token, err := request_app_oath_token(Env_struct)
	if err != nil{
		return "", err
	}
	//twitch_get_user_endpoint := "https://api.twitch.tv/helix/users?"
	bearer_token := "Bearer " + app_oath_token
	url_quary := url.Values{}
	url_quary.Set("id", streamer_id)
	url_encoded_string := url_quary.Encode()
	get_call :=  Env_struct.USER_INFO_URI + url_encoded_string
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

func request_app_oath_token(Env_struct logger.Env_variables) (string, error) {
	//twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
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

func Refresh_token(refresh_token string, user_id string, Env_struct logger.Env_variables) (bool, error){
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
	err = sqlite.Update_tokens(refresh_token_response.Access_token, refresh_token_response.Refresh_token, user_id)
	if err!=nil{
		return false, err
	}
	return true, nil
}

func Test_request_user_oath_token(user_id string, Env_struct logger.Env_variables) (error) {
	//twitch_oath_url := "https://id.twitch.tv/oauth2/token?"
	twitch_oath_url := "http://localhost:8080/auth/authorize"
	url_quary := url.Values{}
	url_quary.Set("client_id", Env_struct.APP_ID)
	url_quary.Set("client_secret", Env_struct.APP_SECRET)
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

func Create_EventSub(user_id string, sub_type string, Env_struct logger.Env_variables)(error){
	client := http.Client{}

	bearer_string := "Bearer " + Env_struct.APP_SECRET

	req_body, err := create_sub_request_body(user_id, sub_type, Env_struct)

	if err!=nil{
		return err
	}

	//Subscription URI info is either "https://api.twitch.tv/helix/eventsub/subscriptions" or "http://localhost:8080/eventsub/subscriptions"

	req, err := http.NewRequest("POST", Env_struct.SUBSCRIPTION_INFO_URI, bytes.NewBuffer(req_body))

	if err!=nil{
		err = errors.New("FILE: twitch_api FUNC: Create_EventSub CALL: http.NewRequest" + err.Error())
		return err
	}

	req.Header.Set("Authorization", bearer_string)
	req.Header.Set("Client-Id", Env_struct.APP_ID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err!=nil || resp.StatusCode != http.StatusOK{
		err = errors.New("FILE: twitch_api FUNC: Create_EventSub CALL: client.Do" + err.Error())
		return err
	}
	defer resp.Body.Close()

	return nil
}

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
			Secret: Env_struct.APP_SECRET},
		}

	body_byte, err := json.Marshal(&body)

	if err != nil{
		err = errors.New("FILE: twitch_api FUNC: create_sub_request_body CALL: json.Marshal " + err.Error())
		return nil, err
	}

	return body_byte, nil
}