package twitch_api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
)

const uri = "http://localhost:3000"

func Generate_state() (string, error) {
	randomBytes := make([]byte, 32)
 
	_, err := rand.Read(randomBytes)

	if err != nil {
		err = errors.New("there was a problem reading the randomBytes byte")
	}

	random_string := base64.RawURLEncoding.EncodeToString(randomBytes)
	fmt.Println(random_string)
	return random_string, err
}

func Scope_requests(request string) (string, error){
	var err error 
	scope := url.Values{}

	if request == "prediction"{
		scope.Set("scope", "channel:manage:predictions")
	} else {
		err = errors.New("Invalid scope request")
	}

	if err != nil {
		return "", err
	}

	url_scope := scope.Encode()
	fmt.Println(url_scope)
	return url_scope, err
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
	
	state, err := Generate_state()

	if err != nil{
		return "", err
	}

	url_quary.Set("state", state)

	encoded_url_quary := url_quary.Encode()

	uri_url_quary := "&redirect_uri=" + uri + "&"

	return_url := url_authorize + uri_url_quary + encoded_url_quary

	fmt.Println(return_url)

	return return_url, err
}