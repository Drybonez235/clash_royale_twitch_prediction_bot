package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	//app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	sqlite "github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/ncruces/go-sqlite3"
)

type Start_royale_bets_json struct {
	Session_id int `json:"session_id"`
	Screen_name string `json:"screen_name"`
	Streamer_player_tag string `json:"streamer_player_tag"`
	Last_refresh_time int `json:"last_refresh_time"`
	Total_points int `json:"total_points"`
}

type Royale_bets_response struct{
	Streamer_info sqlite.Royale_bets_streamer `json:"streamer_info"`
	Leaderboard []sqlite.Leader_board_entry `json:"leaderboard"`
	Battle_results []sqlite.Battle_result `json:"battle_results"`
}

func Start_royale_bets(w http.ResponseWriter, req *http.Request, Env_struct logger.Env_variables, db *sqlite3.Conn)(error){	

	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	

	
	var viewer_json Start_royale_bets_json

	req_body, err := io.ReadAll(req.Body)

	if err!=nil{
		return err
	}

	if err := json.Unmarshal(req_body, &viewer_json); err != nil{
		return err
	}
	
	var viewer sqlite.Royale_bets_viewer
	viewer.Session_id = viewer_json.Session_id
	viewer.Screen_name = viewer_json.Screen_name
	viewer.Streamer_player_tag = viewer_json.Streamer_player_tag
	viewer.Last_refresh_time = viewer_json.Last_refresh_time
	viewer.Total_points = 5000

	streamer_info, leader_board, err := app.Register_viewer(viewer, db)
	if err!=nil{
		return err
	}

	//This should always at least return the one viewer
	if leader_board == nil {
		return errors.New("FILE: royale_bets_server FUNC: Start_royale_bets CALL: app.Register_viwer MESSAGE: Leader Board returned nil")
	}
	//This should always return the streamer that the viewer picked.
	if streamer_info == nil{
		return errors.New("FILE: royale_bets_server FUNC: Start_royale_bets CALL: app.Register_viwer MESSAGE: Streamer info returned nil")
	}

	_, battle_results, err := app.Update_viewer(viewer, Env_struct, db)

	if err!=nil{
		return err
	}

	var response Royale_bets_response

	response.Streamer_info = *streamer_info
	response.Leaderboard = *leader_board
	response.Battle_results = *battle_results

	w.Header().Set("Content-Type:", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err!=nil{
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

//This is what will be called on all subsiquent calls.
func Update_royale_bets(w http.ResponseWriter, req *http.Request, Env_struct logger.Env_variables, db *sqlite3.Conn)(error){
	
	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	var viewer_json Start_royale_bets_json

	req_body, err := io.ReadAll(req.Body)

	if err!=nil{
		return err
	}

	if err := json.Unmarshal(req_body, &viewer_json); err != nil{
		return err
	}
	//UPDATE VIEWER DOES NOT UPDATE THE VIEWER DB.
	var viewer sqlite.Royale_bets_viewer
	viewer.Session_id = viewer_json.Session_id
	viewer.Screen_name = viewer_json.Screen_name
	viewer.Streamer_player_tag = viewer_json.Streamer_player_tag
	viewer.Last_refresh_time = viewer_json.Last_refresh_time
	viewer.Total_points = viewer_json.Total_points



	streamer_info, battle_results, err := app.Update_viewer(viewer, Env_struct, db)
	if err!=nil{
		return err
	}

	if streamer_info == nil{
		return errors.New("FILE: royale_bets_server FUNC: Update_royale_bets CALL: app.Update_viewer MESSAGE: Streamer info returned nil")
	}

	if battle_results == nil{
		return errors.New("FILE: royale_bets_server FUNC: Update_royale_bets CALL: app.Update_viewer MESSAGE: Battles results info returned nil")
	}

	top_ten, err := app.Get_top_ten(viewer, db)

	if err!=nil{
		return err
	}

	var response Royale_bets_response

	response.Streamer_info = *streamer_info
	response.Battle_results = *battle_results
	response.Leaderboard = *top_ten
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err!=nil{
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}

	return nil
}