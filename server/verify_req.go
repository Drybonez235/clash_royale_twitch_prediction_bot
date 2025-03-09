package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
)

func verify_event_message(req *http.Request, body []byte, Env_struct logger.Env_variables)(bool, error){

	fmt.Println("Verify Event Message Fired")
	crafted_message, err := craft_message(req, body)

	if err!=nil{
		return false, err
	}

	created_hmac := create_hmac_message(crafted_message, Env_struct)

	received_message := req.Header.Get("Twitch-Eventsub-Message-Signature")

	if received_message == "" {
		return false, fmt.Errorf("missing Twitch-Eventsub-Message-Signature header")
	}

	parts := strings.SplitN(received_message, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return false, fmt.Errorf("invalid signature format")
	}
	received_hmac, err := hex.DecodeString(parts[1])
	if err != nil {
		fmt.Println("There was some kind of problem")
		return false, fmt.Errorf("error decoding Twitch signature: %v", err)
	}

	if hmac.Equal(created_hmac, received_hmac) {
		return true, nil
	}

	return false, nil
}

func craft_message(req *http.Request, body []byte) (string, error){
	fmt.Println("Crafted Message fired")
	header_string := req.Header.Get("Twitch-Eventsub-Message-Id")

	header_timestamp := req.Header.Get("Twitch-Eventsub-Message-Timestamp")
	if header_string == "" || header_timestamp == ""{
		return "", fmt.Errorf("missing required headers")
	}
	
	return header_string + header_timestamp + string(body), nil
}

func create_hmac_message(crafted_message string, Env_struct logger.Env_variables)( []byte){
	h := hmac.New(sha256.New, []byte(Env_struct.ENCRYPTION_SECRET))
	h.Write([]byte(crafted_message))
	return h.Sum(nil)
}