package clash_royale_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"errors"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
)

const clash_api_url = "https://api.clashroyale.com/v1/players/"

func Get_prior_battles(player_tag string, Env_struct logger.Env_variables)(Match_25, error){
	fmt.Println("Get prior battles fired")
	var battle_log Match_25

	player_tag_url := clash_api_url+"%23"+player_tag+"/battlelog"

	client := http.Client{}

	req, err := http.NewRequest("GET", player_tag_url, nil)

	if err!=nil{
		err = errors.New("FILE: clash_royale_api FUNC: Get_prior_battles CALL: http.NewRequest " + err.Error())
		return battle_log, err
	}

	bearer := "Bearer " + Env_struct.CLASH_API_SECRET

	req.Header.Set("Authorization", bearer)

	resp, err := client.Do(req)

	if err!=nil{
		err = errors.New("FILE: clash_royale_api FUNC: Get_prior_battles CALL: client.Do " + err.Error())
		return battle_log, err
	}

	if resp.StatusCode != http.StatusOK{
		err = errors.New("FILE: clash_royale_api FUNC: Get_prior_battles CALL: client.Do " + resp.Status)
		return battle_log, err
	}
	body, err := io.ReadAll(resp.Body)

	if err!=nil{
		err = errors.New("FILE: clash_royale_api FUNC: Get_prior_battles CALL: io.ReadAll " + err.Error())
		return battle_log, err
	}

	err = json.Unmarshal(body, &battle_log.Matches)
	
	if err!=nil{
		err = errors.New("FILE: clash_royale_api FUNC: Get_prior_battles CALL: json.Unmarshal " + err.Error())
		return battle_log, err
	}
	return battle_log, err
}

func Validate_user_id(player_tag string, Env_struct logger.Env_variables)(bool, error){
	validate_user_url := clash_api_url+"%23"+ player_tag +"/upcomingchests"

	client := http.Client{}

	req, err := http.NewRequest("GET", validate_user_url, nil)

	if err!=nil{
		err = errors.New("FILE: clash_royale_api FUNC: Validate_user_id CALL: http.NewRequest " + err.Error())
		return false, err
	}

	bearer := "Bearer " + Env_struct.CLASH_API_SECRET

	req.Header.Set("Authorization", bearer)

	resp, err := client.Do(req)

	if err!=nil {
		err = errors.New("FILE: clash_royale_api FUNC: Valide_user_id CALL: client.Do " + err.Error())
		return false, err
	}

	if resp.StatusCode != http.StatusOK{
		err = errors.New("FILE: clash_royale_api FUNC: Validate_user_id CALL: client.Do " + resp.Status)
		return false, err
	}

	return true, nil
}