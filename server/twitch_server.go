package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
)

type Authorization_JSON struct {
	code string
	scope string
	state string
}

func Start_server(logger *logger.StandardLogger) {

	fmt.Println("Started server on localhost 3000")

	ticker :=time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C{
		verfify_tokens()
	}

	redirect_uri := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Recieved an app request")
		logger.Info("Recived an app reques")
		proccess_authorization_form(req)
	}

	subscription_callback := func(w http.ResponseWriter, req *http.Request){
		fmt.Println("Subscription call back fired")
		handled, err  := subscription_handler(req, w)

		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError) 
		}

		if !handled {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	handle_event := func(w http.ResponseWriter, req *http.Request){
		err := event_handler(w, req)
		if err!=nil{
			logger.Error(err.Error())
		}
	}

	alive := func(w http.ResponseWriter, req *http.Request){
		fmt.Println("We are here!")
	}

	http.HandleFunc("/redirect", redirect_uri)
	http.HandleFunc("/subscription_handler", subscription_callback)
	http.HandleFunc("/receive_twitch_event", handle_event)
	http.HandleFunc("/", alive)

	

	http.ListenAndServe("localhost:3000", nil)
}

func subscription_handler(req *http.Request, w http.ResponseWriter)(bool,error){
	fmt.Println("Subscription handler function below fired")

	body, err := io.ReadAll(req.Body)
	if err!=nil{
		return false, err
	}
	valid, err := verify_event_message(req, body)
	if err!=nil{
		return false,err
	}
	if !valid{
		return false, nil
	}
	respond_challenge(w, body)
	return true, nil
}


func proccess_authorization_form(req *http.Request)(error){
	fmt.Println("Proccessed auth form")
	
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

	//err = Request_user_oath_token(response.code)

	if err!=nil{
		fmt.Println(err)
		err = errors.New("there was a propbelm with oauth token request")
	}

   return err
}

func event_handler(w http.ResponseWriter, req *http.Request)(error){
	fmt.Println("Recieved an event")
	err := Handle_event(w, req)

	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError) 
		return err	
	}

	return nil
}

func verfify_tokens()(){
	err := Validate_all_tokens()
	if err!=nil{
		fmt.Println(err)
	}
}
 