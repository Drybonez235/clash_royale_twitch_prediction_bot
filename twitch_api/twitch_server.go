package twitch_api

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

type Authorization_JSON struct {
	code string
	scope string
	state string
}

const client_id = "now6dwkymg4vo236ius5d0sn82v9ul"

func Start_server() {

	redirect_uri := func(w http.ResponseWriter, req *http.Request) {
		proccess_authorization_form(req)
	}

	http.HandleFunc("/", redirect_uri)
	http.ListenAndServe("localhost:3000", nil)
}

func proccess_authorization_form(req *http.Request) (Authorization_JSON, error){
	var response Authorization_JSON

	// if req.Response.StatusCode != 200 {
	// 	err := errors.New("the client denied you access to their account")
	// 	return Authorization_JSON{}, err
	// }

	err := req.ParseForm()

	if err != nil{
		err = errors.New("problem reading form values")
		return Authorization_JSON{}, err
	}

	response.scope = req.FormValue("scope")
	response.state = req.FormValue("state")
	response.code = req.FormValue("code")

	valid, err := sqlite.Check_state_nonce(response.state, "state")

	if !valid {
		err = errors.New("invalid state. Malicious request")
	}

	fmt.Println(response)

	err = request_oath_token(response.code)

	if err!=nil{
		err = errors.New("there was a propbelm with oauth token request")
	}

   return response, err
}
