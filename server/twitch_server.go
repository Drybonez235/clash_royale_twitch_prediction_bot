package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

type Authorization_JSON struct {
	code string
	scope string
	state string
}

type Player_tag struct{
	Player_tag string `json:"clash_id""`
}

func Start_server(logger *logger.StandardLogger, Env_struct logger.Env_variables) {
	logger.Info("Started server on localhost 3000")


	redirect_uri := func(w http.ResponseWriter, req *http.Request) {
		logger.Info("Recived an app request")
		valid, err := process_authorization_form(req, Env_struct)

		if err!=nil{
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if !valid{
			logger.Warn("malicious request: State not found in db")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Info("Registered a user from " + req.URL.String())
	}

	subscription_callback := func(w http.ResponseWriter, req *http.Request){

		handled, err  := subscription_handler(req, w, Env_struct)

		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError) 
			logger.Error(err.Error())
		}

		if !handled {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Warn("Subscription callback was not handled correctly")
			return
		}
	}

	handle_event := func(w http.ResponseWriter, req *http.Request){
		err := event_handler(w, req, logger, Env_struct)
		if err!=nil{
			logger.Error(err.Error())
		}
		logger.Info("Logged event" + req.URL.String())
	}

	verify_player_tag := func(w http.ResponseWriter, req *http.Request){
		err := handle_verify_player_tag(w, req, Env_struct)
		if err!=nil{
			logger.Error(err.Error())
		}
	}


	http.HandleFunc("/redirect", redirect_uri)
	http.HandleFunc("/subscription_handler", subscription_callback)
	http.HandleFunc("/receive_twitch_event", handle_event)
	http.HandleFunc("/verify_player_tag", verify_player_tag)
	

	http.ListenAndServe("localhost:3000", nil)
}

func subscription_handler(req *http.Request, w http.ResponseWriter, Env_struct logger.Env_variables)(bool,error){

	body, err := io.ReadAll(req.Body)
	if err!=nil{
		return false, err
	}
	valid, err := verify_event_message(req, body)
	if err!=nil{
		return false, err
	}
	if !valid{
		return false, err
	}
	respond_challenge(w, body)
	return true, nil
}


func process_authorization_form(req *http.Request, Env_struct logger.Env_variables)(bool, error){
	fmt.Println("Proccessed auth form")
	
	var response Authorization_JSON

	err := req.ParseForm()

	if err != nil{
		err = errors.New("problem reading authorization form values")
		return false, err
	}

	response.scope = req.FormValue("scope")
	response.state = req.FormValue("state")
	response.code = req.FormValue("code")

	valid, player_tag ,err := sqlite.Check_state_nonce(response.state, "state")

	if err !=nil {
		return false, err
	}

	if !valid{
		return false, nil
	}

	if player_tag == ""{
		err = errors.New("player tag associated with state was blank")
		return false,err
	}

	err = twitch_api.Request_user_oath_token(response.code, player_tag)

	if err!=nil{
		return false, err
	}

   return true, nil
}

func event_handler(w http.ResponseWriter, req *http.Request, logger *logger.StandardLogger, Env_struct logger.Env_variables)(error){
	err := Handle_event(w, req, logger)

	if err!=nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err	
	}
	return nil
}

func verfify_tokens()(error){
	err := Validate_all_tokens()
	if err!=nil{
		return err
	}
	return nil
}

func handle_verify_player_tag(w http.ResponseWriter, req *http.Request, Env_struct logger.Env_variables)(error){

	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	body, err := io.ReadAll(req.Body)
	if err!=nil{
		return err}

	var string_player_tag Player_tag

	err = json.Unmarshal(body, &string_player_tag)

	if err!=nil{
		return err
	}

	true_false, err := clash.Validate_user_id(string_player_tag.Player_tag)

	if err!=nil{return nil}

	if true_false{
		authorize_app_url, nonce, err := twitch_api.Generate_authorize_app_url(App_id, "prediction")
		if err!=nil{return err}

		err = sqlite.Update_nonce(nonce, string_player_tag.Player_tag)

		if err!=nil{return err}

		json_string := fmt.Sprintf(`{"valid":true,"URL":"%s"}`, authorize_app_url)
		w.Write([]byte(json_string))
	} else {	
		w.Write([]byte(`{"valid":false}`))
	}
	return nil
}

 