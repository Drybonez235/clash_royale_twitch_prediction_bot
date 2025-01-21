package twitch_api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

const uri = "http://localhost:3000"

func Generate_state_nonce(state_nonce string) ( string, error) {
	randomBytes := make([]byte, 32)
 
	_, err := rand.Read(randomBytes)

	if err != nil {
		err = errors.New("there was a problem reading the randomBytes byte")
		return "", err
	}

	random_string := base64.RawURLEncoding.EncodeToString(randomBytes)
	fmt.Println(random_string)

	table := ""

	if state_nonce == "state"{
		table = "state"
	} else if state_nonce == "nonce" {
		table = "nonce"
	} else {
		err = errors.New("invalid table given")
	}
	if err != nil{
		return "", err
	}

	err = sqlite.Write_state_nonce(random_string, table)

	if err!= nil{
		return "", err
	}

	return random_string, err
}

func Scope_requests(request string) (string, error){
	var err error 

	request_string := ""

	if request == "prediction"{
		request_string = "channel:manage:predictions openid" 
		return request_string, err
	} else {
		err = errors.New("invalid scope request")
	}

	if err != nil {
		return "", err
	}

	fmt.Println(request)
	return request, err
}

func Generate_authorize_app_url(client_id string, scope_request string)(string, error){
	
	url_authorize := "https://id.twitch.tv/oauth2/authorize?"

	url_quary := url.Values{}

	url_quary.Set("client_id", client_id)
	url_quary.Set("force_verify", "false")
	url_quary.Set("response_type", "code")

	scope, err := Scope_requests(scope_request)

	if err != nil{
		return "", err
	}

	url_quary.Set("scope", scope)
	
	state, err := Generate_state_nonce("state")

	if err != nil{
		return "", err
	}


	url_quary.Set("state", state)

	nonce, err := Generate_state_nonce("nonce")

	if err != nil{
		return "", err
	}

	url_quary.Set("nonce", nonce)

	if err != nil{
		return "", err
	}
	encoded_url_quary := url_quary.Encode()

	uri_url_quary := "&redirect_uri=" + uri + "&"

	return_url := url_authorize + uri_url_quary + encoded_url_quary

	fmt.Println(return_url)

	return return_url, err
}