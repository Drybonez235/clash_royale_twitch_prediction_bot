package app

import (
	"errors"
	"time"
	"fmt"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/ncruces/go-sqlite3"
)

//Register_viewer is the first call the website makes. It gets the streamer info and sets it if it has not been set.
//It returns the top ten viewers for the leaderboard and the streamer info.
func Register_viewer(viewer sqlite.Royale_bets_viewer, db *sqlite3.Conn) (*sqlite.Royale_bets_streamer, *[]sqlite.Leader_board_entry, error) {
	if err := sqlite.Insert_royale_bets_viewer(db, viewer); err !=nil{
		err = errors.New("FILE: royale_bets_app FUNC: Register_Viewer CALL: sqlite.Insert_Royale_bets_viewer " + err.Error())
		return nil, nil, err
	}

	streamer, err := sqlite.Get_royale_bets_streamer(db, viewer.Streamer_player_tag, viewer.Session_id)

	if err!=nil{
		return  nil, nil,errors.New("FILE: royale_bets_app FUNC: Register_viewer CALL: sqlite.Get_royale_bets_streamer " + err.Error())
	}

	if streamer == nil{
		var new_streamer sqlite.Royale_bets_streamer

		// time := int(time.Now().Unix())

		new_streamer.Losses = 0
		new_streamer.Wins = 0
		new_streamer.Streamer_player_tag = viewer.Streamer_player_tag
		new_streamer.Stream_start_time = viewer.Session_id 
		new_streamer.Streamer_last_refresh_time = viewer.Session_id

		if err = sqlite.Insert_royale_bets_streamer(db, new_streamer); err!=nil{
			return nil, nil ,errors.New("FILE: royale_bets_app FUNC: Register_viewer CALL: sqlite.Insert_royale_bets_streamer " + err.Error())
		}
		streamer = &new_streamer
	}

	top_ten, err := Get_top_ten(viewer, db)

	if err!=nil{
		return nil, nil, errors.New("FILE: royale_bets_app FUNC: Register_viewer CALL: Get_top_ten " + err.Error())
	} 

	return streamer, top_ten, nil
}

//This is the subsiquent call the website makes to update. It gets streamer info, checks to see if there is an update availible. If there is an update it updates. 
func Update_viewer(viewer sqlite.Royale_bets_viewer, Env_struct logger.Env_variables, db *sqlite3.Conn)(*sqlite.Royale_bets_streamer, *[]sqlite.Battle_result, error){
	
	err := sqlite.Update_royale_bets_viewer(db, viewer.Session_id, viewer.Screen_name, viewer.Total_points, viewer.Last_refresh_time)

	if err!=nil{
		return nil, nil, errors.New("FILE: royale_bets_app FUNC: sqlite.Update_royale_bets_viewer " + err.Error())
	}

	streamer_info, err := sqlite.Get_royale_bets_streamer(db, viewer.Streamer_player_tag, viewer.Session_id)
	if err!=nil{
		return nil, nil, errors.New("FILE: royale_bets_app FUNC: Update_viewer CALL: sqlite.Get_royale_bets_streamer " + err.Error())
	}

	//If the streamer was updated before the viewer last updated, then there could be a new battle.
	if streamer_info.Streamer_last_refresh_time <= viewer.Last_refresh_time{
		//This calls the clash royale api.
		if err := Update_streamer_battles(streamer_info.Streamer_player_tag, viewer.Session_id ,Env_struct, db); err!=nil{
			return streamer_info, nil, err
		}

	//If the streamer was updated after the viewer, then there could be a battle in the database already.
	} else if streamer_info.Streamer_last_refresh_time > viewer.Last_refresh_time{
		//This calls our database
		new_battles, err := sqlite.Get_battle_result(db, streamer_info.Streamer_player_tag, viewer.Last_refresh_time)
		if err!=nil{
			return nil, nil, err
		} 
		return streamer_info, &new_battles, nil
	}

	new_battles, err := sqlite.Get_battle_result(db, streamer_info.Streamer_player_tag, viewer.Last_refresh_time)
		if err!=nil{
			return nil, nil, err
		} 
		// if len(new_battles) == 0{
		// 	//empty_array := []sqlite.Battle_result{}
		// 	//Maybe. Maybe not. I have to decide what to rteturn. Proably in the server layer.
		// 	return streamer_info, &new_battles, nil //errors.New("FILE: royale_bets_app FUNC: Update_viewer CALL: sqlite.Get_battle_result MESSAGE: No new_battles")
		// }
	return streamer_info, &new_battles, nil	
}

//This calls the clash royale api and adds them to the db if the battle time is less  
func Update_streamer_battles(streamer_tag string, viewer_session_id int, Env_struct logger.Env_variables, db *sqlite3.Conn) error {

	matches, err := clash_royale_api.Get_prior_battles(streamer_tag, Env_struct)
	if err != nil{
		return errors.New("FILE: royale_bets_app FUNC: Update_streamer_battles CALL: clash_royale_api.Get_prior_battles " + err.Error())
	}
	//Maybe this could be the culprit? I should check but this should always return the last 25 battles of the already verified clash tag
	if len(matches.Matches) == 0{
		return nil
	}

	//streamer last update time is not updating propperly
	streamer, err := sqlite.Get_royale_bets_streamer(db, streamer_tag, viewer_session_id)

	if err!=nil{
		return errors.New("FILE: royale_bets_app FUNC: Update_streamer_battles CALL: sqlite.Get_royale_bets_streamer " + err.Error())
	}

	if streamer == nil{
		return errors.New("FILE: royale_bets_app FUNC: Update_streamer_battles CALL: sqlite.Get_royale_bets_streamer ERROR: Streamer not found")
	}

	for i:=0; i<len(matches.Matches); i++{
		battle := matches.Matches[i]
		battle_time, err := clash_royale_api.String_time_to_time_time(battle.BattleTime)

		fmt.Println("Battle Time ", int(battle_time.Unix() * 1000), " Streamer Last refresh time ", streamer.Streamer_last_refresh_time)

		if err != nil{
			return errors.New("FILE: royale_bets_app FUNC: Update_streamer_battles CALL: clash_royale_api.String_time_to_time_time " + err.Error())
		}

		if streamer.Streamer_last_refresh_time > (int(battle_time.Unix()) * 1000){
			break
		} else {
			var result sqlite.Battle_result

			result.Battle_time = int(battle_time.Unix()) * 1000
			result.Player_tag = battle.Team[0].Tag[1:]
			result.Red_crowns_taken = battle.Team[0].Crowns
			result.Blue_crowns_lost = battle.Opponent[0].Crowns
	
			if err = sqlite.Insert_battle_result(db, result); err != nil{
				return errors.New("FILE: royale_bets_app FUNC: Update_streamer_battles CALL: sqlite.Insert_battle_result " + err.Error())
			}

			//This is untested but it should increase the win or loss value of the streamer by 1 depedning on the battle result.
			if result.Red_crowns_taken > result.Blue_crowns_lost{
				fmt.Println("More red crowns than Blue crowns")
				if err = sqlite.Update_royale_bets_streamer_wins_losses(db, streamer.Streamer_player_tag, streamer.Stream_start_time, int(time.Now().Unix() * 1000), "win"); err!=nil{
					return err
				}
			} else if result.Blue_crowns_lost > result.Red_crowns_taken{
				fmt.Println("More Blue crowns than red crowns")
				if err = sqlite.Update_royale_bets_streamer_wins_losses(db, streamer.Streamer_player_tag, streamer.Stream_start_time, int(time.Now().Unix() * 1000), "lose"); err!=nil{
					return err
				}	
			} 
		}
	}
	return nil
}
//This returns the top ten and the rank and points of the viewer. There should always be at least one entry.
func Get_top_ten(Viewer sqlite.Royale_bets_viewer, db *sqlite3.Conn) (*[]sqlite.Leader_board_entry, error){
	
	streamer_info, err := sqlite.Get_royale_bets_streamer(db, Viewer.Streamer_player_tag, Viewer.Session_id)

	if err!=nil{ 
		return nil, errors.New("FILE: royale_bets_app FUNC: Get_top_ten CALL: sqlite.Get_royale_bets_streamer " + err.Error())
	}

	if streamer_info == nil{
		return nil, errors.New("FILE: royale_bets_app FUNC: Get_top_ten CALL: sqlite.Get_royale_bets_streamer MESSAGE: Streamer not found")
	}

	entries, err := sqlite.Get_top_ten_and_viewer_position(db, Viewer.Streamer_player_tag, streamer_info.Stream_start_time, Viewer.Session_id)

	if err!=nil{
		return nil, errors.New("FILE: royale_bets_app FUNC: Get_top_ten CALL: sqlite.Get_top_ten_and_viewer_position " + err.Error())
	}

	if len(*entries) == 0 {
		return nil, errors.New("FILE: royale_bets_app FUNC: Get_top_ten CALL: sqlite.Get_top_ten_and_viewer_position MESSAGE: No viewers in viewer table")
	}

	return entries, nil
}