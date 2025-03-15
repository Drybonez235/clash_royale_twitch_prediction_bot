package server

import (
	//"bytes"
	"encoding/json"
	//"errors"
	"net/http"
	//"fmt"
	//logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	// "github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	// "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

func respond_challenge(w http.ResponseWriter, body []byte){
	var data Challenge_struct
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	challenge := data.Challenge
	if challenge ==""{
		http.Error(w, "Challenge not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
}
