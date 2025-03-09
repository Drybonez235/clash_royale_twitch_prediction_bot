package twitch_api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	//"fmt"
	"net/url"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

func Generate_state_nonce(state_nonce string) ( string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		err = errors.New("there was a problem reading the randomBytes byte")
		return "", err
	}
	random_string := base64.RawURLEncoding.EncodeToString(randomBytes)
	table := ""
	if state_nonce == "state"{
		table = "state"
		err = sqlite.Write_state_nonce(random_string, table)
	} else if state_nonce == "nonce" {
		return random_string, err
	} else {
		err = errors.New("invalid table given")
	}
	if err != nil{
		return "", err
	}
	return random_string, err
}

//This function generates the url that streamers will use to connect to Twitch. It returns a URL and a nonce, and an error.
func Generate_authorize_app_url(scope_request string, Env_struct logger.Env_variables)(string, string, error){
	url_quary := url.Values{}
	url_quary.Set("client_id", Env_struct.APP_ID)
	url_quary.Set("force_verify", "false")
	url_quary.Set("response_type", "code")
	scope, err := Scope_requests(scope_request)
	if err != nil{
		return "", "",err
	}
	url_quary.Set("scope", scope)
	state, err := Generate_state_nonce("state")
	if err != nil{
		return "", "", err
	}
	url_quary.Set("state", state)
	nonce, err := Generate_state_nonce("nonce")
	if err != nil{
		return "", "", err
	}
	url_quary.Set("nonce", nonce)
	encoded_url_quary := url_quary.Encode()
	uri_url_quary := "&redirect_uri=" + Env_struct.ROYALE_BETS_URL + "/redirect" + "&"
	return_url := Env_struct.OAUTH_AUTHORIZE_URI + uri_url_quary + encoded_url_quary
	return return_url, nonce, nil
}