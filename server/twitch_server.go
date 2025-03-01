package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

type Authorization_JSON struct {
	code string
	scope string
	state string
}

func Start_server(logger *logger.StandardLogger) {
	logger.Info("Started server on localhost 3000")
	defer logger.Info("Server stopped for some reason")
	ticker :=time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C{
		err := verfify_tokens()
		if err!=nil{
			logger.Error(err.Error())
		}
		logger.Info("Verified tokens")
	}

	redirect_uri := func(w http.ResponseWriter, req *http.Request) {
		logger.Info("Recived an app request")
		valid, err := proccess_authorization_form(req)

		if err!=nil{
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		} 
		
		if !valid{
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
		}

		logger.Info("Registered a user from " + req.URL.String())
	}

	subscription_callback := func(w http.ResponseWriter, req *http.Request){

		handled, err  := subscription_handler(req, w)

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
		err := event_handler(w, req, logger)
		if err!=nil{
			logger.Error(err.Error())
		}
		logger.Info("Logged event" + req.URL.String())
	}


	http.HandleFunc("/redirect", redirect_uri)
	http.HandleFunc("/subscription_handler", subscription_callback)
	http.HandleFunc("/receive_twitch_event", handle_event)
	

	http.ListenAndServe("localhost:3000", nil)
}

func subscription_handler(req *http.Request, w http.ResponseWriter)(bool,error){

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


func proccess_authorization_form(req *http.Request)(bool, error){
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

	valid, err := sqlite.Check_state_nonce(response.state, "state")

	if err !=nil {
		return false, err
	}

	if !valid{
		err = errors.New("malicious request: State not found in db")
		return false, err
	}

	err = twitch_api.Request_user_oath_token(response.code)

	if err!=nil{
		return false, err
	}

   return true, nil
}

func event_handler(w http.ResponseWriter, req *http.Request, logger *logger.StandardLogger)(error){
	err := Handle_event(w, req, logger)

	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError) 
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
 