package clash_royale_api

import (
	"time"
	"fmt"
)

func New_battle(player_tag string, prediction_created_at time.Time)error{
	Matches, err := Get_prior_battles(player_tag)

	if err!=nil{
		return err
	}

	for i:=0; i < len(Matches.Matches); i++{
		Match := Matches.Matches[i]
		Battle_time, err := String_time_to_time_time(Match.BattleTime)

		if err!=nil{
			return err
		}

		if Battle_time.After(prediction_created_at){
			//Check for win or lose
			//End prediction
		} else{
			break
		}
	}


	fmt.Println(Matches)
	//We are going to have to mess with time 
	return nil
}

func String_time_to_time_time(battle_time_string string)(time.Time, error){

	battle_time_string_year := battle_time_string[0:4]
	battle_time_string_month := battle_time_string[4:6]
	battle_time_string_day := battle_time_string[6:8]
	battle_time_string_hour :=  battle_time_string[9:11] 
	battle_time_string_minute := battle_time_string[11:13] 
	battle_time_string_second :=  battle_time_string[13:15]
	
	battle_time_string_string := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ",battle_time_string_year, battle_time_string_month, battle_time_string_day, battle_time_string_hour, battle_time_string_minute, battle_time_string_second)
	
	fmt.Println(battle_time_string_string)

	t, err := time.Parse(time.RFC3339, battle_time_string_string)

	if err!=nil{
		return  t, err
	}

	return t, nil
}