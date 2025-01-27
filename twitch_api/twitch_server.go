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

func Start_server() {

	redirect_uri := func(w http.ResponseWriter, req *http.Request) {
		proccess_authorization_form(req)
	}

	http.HandleFunc("/redirect", redirect_uri)
	http.ListenAndServe("localhost:3000", nil)
}

func proccess_authorization_form(req *http.Request) (error){
	var response Authorization_JSON

	err := req.ParseForm()

	if err != nil{
		err = errors.New("problem reading form values")
		return err
	}

	response.scope = req.FormValue("scope")
	response.state = req.FormValue("state")
	response.code = req.FormValue("code")

	valid, err := sqlite.Check_state_nonce(response.state, "state")

	if !valid || err !=nil {
		err = errors.New("invalid state. Malicious request or the check state didn't work")
		return err
	}

	fmt.Println(response)

	err = Request_user_oath_token(response.code)

	if err!=nil{
		fmt.Println(err)
		err = errors.New("there was a propbelm with oauth token request")
	}

   return err
}
