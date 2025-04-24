package server

import(
	"net/http"
	"io"
	"errors"
	"encoding/json"
	"fmt"
	
	twitch_api "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	sqlite "github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	clash 	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	"github.com/ncruces/go-sqlite3"
	
)

func Handle_verify_player_tag(w http.ResponseWriter, req *http.Request, Env_struct logger.Env_variables, db *sqlite3.Conn)(error){

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
		err = errors.New("FILE twitch_server FUNC: handle_verify_player_tag CALL: io.ReadAll " + err.Error())
		return err}

	var string_player_tag verify_player_tag_req

	err = json.Unmarshal(body, &string_player_tag)

	if err!=nil{
		err = errors.New("FILE twitch_server FUNC: handle_verify_player_tag CALL: json.Unmarshal " + err.Error())
		return err
	}

	Top_25_matched, err := clash.Get_prior_battles(string_player_tag.Player_tag, Env_struct)
	if err!=nil{return nil}

	Match_array := Top_25_matched.Matches

	if len(Match_array) != 0{
		if string_player_tag.Req_page == "streamer_sign_up"{
		authorize_app_url, state, err := twitch_api.Generate_authorize_app_url("prediction", Env_struct, db)
		if err!=nil{return err}

		err = sqlite.Update_state(db ,state, string_player_tag.Player_tag)

		if err!=nil{return err}

		json_string := fmt.Sprintf(`{"valid":true,"URL":"%s"}`, authorize_app_url)
		w.Write([]byte(json_string))

		} else if string_player_tag.Req_page == "player_tag_input"{
			clash_name := Match_array[0].Team[0].Name
			w.Write([]byte(fmt.Sprintf(`{"valid":true, "clash_name": "%s"}`, clash_name)))
		}
	} else {	
		w.Write([]byte(`{"valid":false}`))
	}
	return nil
}

func Event_handler(w http.ResponseWriter, req *http.Request, logger *logger.StandardLogger, Env_struct logger.Env_variables, db *sqlite3.Conn)(error){
	err := Handle_event(w, req, logger, Env_struct, db)

	if err!=nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err	
	}
	return nil
}

func Process_authorization_form(req *http.Request, Env_struct logger.Env_variables, db *sqlite3.Conn)(bool, error){
	var response Authorization_JSON
	err := req.ParseForm()
	if err != nil{
		err = errors.New("FILE: twitch_server FUNC: proccess_authorization_form CALL: req.ParseForm() "+ err.Error())
		return false, err
	}

	response.scope = req.FormValue("scope")
	response.state = req.FormValue("state")
	response.code = req.FormValue("code")

	valid, player_tag, err := sqlite.Check_state_nonce(db, response.state, "state")

	if err !=nil {
		return false, err
	}
	if !valid{
		return false, nil
	}
	if player_tag == ""{
		err = errors.New("FILE twitch_server FUNC: proccess_authorization_form BUG: player_tag was blank")
		return false,err
	}

	err = twitch_api.Request_user_oath_token(response.code, player_tag, Env_struct, db)

	if err!=nil{
		return false, err
	}

   return true, nil
}

func Subscription_handler(req *http.Request, w http.ResponseWriter, Env_struct logger.Env_variables)(bool,error){

	body, err := io.ReadAll(req.Body)
	if err!=nil{
		err = errors.New("FILE twitch_server FUNC: subscription_hadnler CALL: io.ReadALl " + err.Error())
		return false, err
	}
	valid, err := verify_event_message(req, body, Env_struct)
	if err!=nil{
		return false, err
	}
	if !valid{
		return false, nil
	}
	respond_challenge(w, body)
	return true, nil
}