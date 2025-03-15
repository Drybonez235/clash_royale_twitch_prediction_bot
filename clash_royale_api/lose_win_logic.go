package clash_royale_api

import (
	"fmt"
	"time"
	"errors"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
)

func New_battle(player_tag string, prediction_created_at time.Time, Env_struct logger.Env_variables)(string, error){
	Matches, err := Get_prior_battles(player_tag, Env_struct)
	if err!=nil{
		return "error", err
	}

	for i:=0; i < len(Matches.Matches); i++{
		Match := Matches.Matches[i]
		Battle_time, err := String_time_to_time_time(Match.BattleTime)

		if err!=nil{
			return "", err
		}

		if Battle_time.After(prediction_created_at){
			lose_win, err := Lose_win(Match)

			if err!=nil{
				return "",err
			}

			return lose_win, nil
		} else{
			return "no_new_battles", nil
		}
	}
	return "err", nil
}

func String_time_to_time_time(battle_time_string string)(time.Time, error){
	battle_time_string_year := battle_time_string[0:4]
	battle_time_string_month := battle_time_string[4:6]
	battle_time_string_day := battle_time_string[6:8]
	battle_time_string_hour :=  battle_time_string[9:11] 
	battle_time_string_minute := battle_time_string[11:13] 
	battle_time_string_second :=  battle_time_string[13:15]
	
	battle_time_string_string := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ",battle_time_string_year, battle_time_string_month, battle_time_string_day, battle_time_string_hour, battle_time_string_minute, battle_time_string_second)
	t, err := time.Parse(time.RFC3339, battle_time_string_string)
	if err!=nil{
		err = errors.New("FILE: lose_win_logic FUNC: String_time_to_time_time CALL: time.Parse " + err.Error())
		return  t, err
	}

	return t, nil
}

func Lose_win(match Match) (string, error) {
	if len(match.Team) == 0 {
		return "", errors.New("FILE: lose_win_logic FUNC: Lose_win BUG: match.Team is empty")
	}
	if len(match.Opponent) == 0 {
		return "", errors.New("FILE: lose_win_logic FUNC: Lose_win BUG: match.Opponent is empty")
	}

	streamer_player := match.Team[0]
	opponent_player := match.Opponent[0]

	streamer_player_crowns := streamer_player.Crowns
	opponent_player_crowns := opponent_player.Crowns

	if streamer_player_crowns > opponent_player_crowns {
		return "win", nil
	} else if opponent_player_crowns > streamer_player_crowns {
		return "lose", nil
	} else {
		return "tie", nil
	}
}