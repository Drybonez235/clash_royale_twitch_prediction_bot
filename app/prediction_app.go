package app

import (
	"fmt"
	"time"
	"errors"
	clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	sqlite "github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

func Start_prediction_app(sub string, Env_struct logger.Env_variables) error {
	fmt.Println("Started Prediction App")
	user, err := sqlite.Get_twitch_user("sub", sub)
	if err!=nil{
		return err
	}

	stream, err := sqlite.Twitch_user_online(sub)

	if err!=nil{
		return err
	}


	for stream {
		fmt.Println("We did start the stream loop")
		//We have to check to see if there is an active prediction here that was set by me. IF it was not set by me, then we need to wait.
		prediction_id, _, err := sqlite.Get_predictions(user.User_id, "ACTIVE")

		if prediction_id == "null"{
			prediction_id = ""
		}

		if err!=nil{return err}
		own_active_prediction, err := twitch.Check_prediction(sub, user.Access_token, prediction_id, Env_struct)
		if err!= nil{
			return err
		}

		 if own_active_prediction == "our_prediction" {
			fmt.Println("Our prediction")
			err = Watch_prediction(sub, user, Env_struct)
			if err!=nil{
				return err
			}

		 }else if own_active_prediction == "no_active_prediction"{
			fmt.Println("No active prediction")
			err = twitch.Start_prediction(user, Env_struct)
			if err!=nil{
				return err
			}
		} else if own_active_prediction == "not_our_prediction" {
			fmt.Println("Not our prediction")
			time.Sleep(30 * time.Second)
		} else {
			err = errors.New("FILE: prediction_app FUNC: Start_prediction_app BUG: own_active_prediction invalid")
			return err
		} 

		stream, err = sqlite.Twitch_user_online(sub)

		if err!=nil{
			stream = false
			return err
		}
	}
	return nil
}

//IF we own the prediction then this function will fire.
func Watch_prediction(sub string, user sqlite.Twitch_user, Env_struct logger.Env_variables)error{
	prediction_id, created_at, err := sqlite.Get_predictions(sub, "ACTIVE")
	if err!=nil{return err}
	if prediction_id == ""{
		err = errors.New("FILE: prediction_app FUNC: Watch_prediction BUG: prediction_id was blank")
		return err
	}

	t_created_at, err := time.Parse(time.RFC3339, created_at)

	if err!=nil{
		err = errors.New("FILE: prediction_app FUNC Watch_prediction CALL: time.Parse " + err.Error())
		return err
	}
	new_battle := "no_new_battles"
	for (new_battle == "no_new_battles"){
		fmt.Println("No new battles")
		user, err := sqlite.Get_twitch_user("sub", sub)
		if err!=nil{return err}
		if user.Online != 1{
			return nil
		}

		new_battle, err := clash.New_battle(user.Player_tag, t_created_at, Env_struct)

		if err!=nil{
			return err
		}

		if new_battle == "no_new_battles"{
			time.Sleep(30 * time.Second)
		} else if new_battle == "win"{
			fmt.Println("You won")
			outcome_id, err := sqlite.Get_prediction_outcome_id(prediction_id, 1)
			if err!=nil{
				return err
			}
			if outcome_id == ""{
				err = errors.New("FILE: prediction_app FUNC: Watch_prediction BUG: outcome_id was blank")
				return err
			}
			err = twitch.End_prediction(prediction_id, outcome_id, user.User_id, user.Access_token, Env_struct)
			if err !=nil{
				return err
			}
			return nil
		} else if new_battle == "lose"{
			fmt.Println("You lost")
			outcome_id, err := sqlite.Get_prediction_outcome_id(prediction_id, 0)
			if err!=nil{
				return err
			}
			if outcome_id == ""{
				err = errors.New("FILE: prediction_app FUNC: Watch_prediction BUG: outcome_id was blank")
				return err
			}
			err = twitch.End_prediction(prediction_id, outcome_id, user.User_id, user.Access_token, Env_struct)
			if err !=nil{
				return err
			}
		return nil
	} else if new_battle == "tie"{
		fmt.Println("Battle tie")
	}
}
	return nil
}