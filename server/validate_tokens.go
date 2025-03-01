package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

type Refresh_token_response struct{
	Access_token string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Scope []string `json:"scope"`
	Token_type string `json:"token_type"`
}

//This iterates through all current access tokens and verifys their validity. If it isn't valid, it attempts to refresh using the refresh token. If the refresh fails, then it deletes the user from the db. (Maybe in the future it could email)
func Validate_all_tokens()(error){
	var Token_list []sqlite.Twitch_user_refresh

	Token_list, err := sqlite.Get_all_access_tokens()

	for i:=0; i<len(Token_list); i++{
		Twitch_user := Token_list[i]
		
		valid, err := Validate_token(Twitch_user)

		if err!=nil{return err}

		if !valid{
			fmt.Println("Deleted a user due to invalid refresh attempt")
		}
	}

	return err
}

func Validate_token(Twitch_user sqlite.Twitch_user_refresh)(bool, error){
	fmt.Println("Validate token ran")

	twitch_validation_endpoint := "https://id.twitch.tv/oauth2/validate"

	client := &http.Client{}

	req, err := http.NewRequest("GET", twitch_validation_endpoint, nil)
	if err != nil{
		err = errors.New("there was something wrong with the GET request")
		return false, err	
	}
	req.Header.Set("Authorization", "OAuth " + Twitch_user.Access_token )
	resp, err := client.Do(req)
	if err != nil{
		err = errors.New("there was something wrong with the GET response")
		return false, err	
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		return true, nil

	} else if resp.StatusCode == 401 {

		refreshed, err := Refresh_token(Twitch_user)
		if err !=nil{
			return false, err
		}

		if !refreshed {
			err = sqlite.Remove_twitch_user(Twitch_user.User_id)
			if err!=nil{
				return false, err
			}
		}

		return true, nil
	}
	return true, err
}
func Refresh_token(Twitch_user sqlite.Twitch_user_refresh) (bool, error){
	refresh_token_url := "https://id.twitch.tv/oauth2/token"
	client := &http.Client{}
	url_quary := url.Values{}
	url_quary.Set("client_id", App_id)
	url_quary.Set("client_secret", app_secret)
	url_quary.Set("grant_type", "refresh_token")
	url_quary.Set("refresh_token", Twitch_user.Refresh_token)
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
	err = sqlite.Update_tokens(refresh_token_response.Access_token, refresh_token_response.Refresh_token, Twitch_user.User_id)
	if err!=nil{
		return false, err
	}
	return true, nil
}