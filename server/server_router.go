package server

import (
	"net/http"

	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/ncruces/go-sqlite3"
)

type Authorization_JSON struct {
	code string
	scope string
	state string
}

type verify_player_tag_req struct{
	Player_tag string `json:"clash_id"`
	Req_page string `json:"req_page"`
}

func Start_server(logger *logger.StandardLogger, Env_struct logger.Env_variables, db *sqlite3.Conn) {
	logger.Info("Started server on localhost 3000")

	redirect_uri := func(w http.ResponseWriter, req *http.Request) {
		logger.Info("Recived an app request")
		valid, err := Process_authorization_form(req, Env_struct, db)

		if err!=nil{
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if !valid{
			logger.Warn("malicious request: State not found in db")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.Redirect(w, req, Env_struct.ROYALE_BETS_REDIRECT_URL, http.StatusSeeOther)
		logger.Info("Registered a user from " + req.URL.String())
	}

	subscription_callback := func(w http.ResponseWriter, req *http.Request){

		handled, err  := Subscription_handler(req, w, Env_struct)

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
		err := Event_handler(w, req, logger, Env_struct, db)
		if err!=nil{
			logger.Error(err.Error())
		}
		logger.Info("Logged event" + req.URL.String())
	}

	verify_player_tag := func(w http.ResponseWriter, req *http.Request){
		err := Handle_verify_player_tag(w, req, Env_struct, db)
		if err!=nil{
			logger.Error(err.Error())
		}
	}

	start_royale_bets := func(w http.ResponseWriter, req *http.Request){
		err := Start_royale_bets(w, req, Env_struct, db)
		if err!=nil{
			logger.Error(err.Error())
		}
	}

	update_royale_bets := func(w http.ResponseWriter, req *http.Request){
		err := Update_royale_bets(w, req, Env_struct, db)
		if err!=nil{
			logger.Error(err.Error())
		}
	}

	http.HandleFunc("/redirect", redirect_uri)
	http.HandleFunc("/subscription_handler", subscription_callback)
	http.HandleFunc("/receive_twitch_event", handle_event)
	http.HandleFunc("/verify_player_tag", verify_player_tag)
	http.HandleFunc("/start_royale_bets", start_royale_bets)
	http.HandleFunc("/update_royale_bets", update_royale_bets)
	
	http.ListenAndServe("localhost:3000", nil)
}
 