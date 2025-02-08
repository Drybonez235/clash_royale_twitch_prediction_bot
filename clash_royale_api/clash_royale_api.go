package clash_royale_api

import (
	"encoding/json"
	"io"
	"net/http"
	//"fmt"
)

const access_token = ""
const clash_api_url = "https://api.clashroyale.com/v1/players/"

func Get_prior_battles(player_tag string)(Match_25, error){
	var battle_log Match_25

	player_tag_url := clash_api_url+"%23"+player_tag+"/battlelog"

	client := http.Client{}

	req, err := http.NewRequest("GET", player_tag_url, nil)

	if err!=nil{
		return battle_log, err
	}

	bearer := "Bearer " + access_token

	req.Header.Set("Authorization", bearer)

	resp, err := client.Do(req)

	if err!=nil || resp.StatusCode!= http.StatusOK{
		return battle_log, err
	}

	body, err := io.ReadAll(resp.Body)

	if err!=nil{
		return battle_log, err
	}

	err = json.Unmarshal(body, &battle_log.Matches)

	return battle_log, err
}

func Validate_user_id(player_tag string)(bool, error){
	validate_user_url := clash_api_url+"%23"+ player_tag +"/upcomingchests"

	client := http.Client{}

	req, err := http.NewRequest("GET", validate_user_url, nil)

	if err!=nil{
		return false, err
	}

	bearer := "Bearer " + access_token

	req.Header.Set("Authorization", bearer)

	resp, err := client.Do(req)

	if err!=nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK{
		return false, err
	}

	return true, nil
}