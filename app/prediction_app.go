package app

import (
	"fmt"
	"time"

	clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	sqlite "github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

func Start_prediction_app(sub string) error {
	fmt.Println("Started the prediction app.")
	
	user, err := sqlite.Get_twitch_user("sub", sub)
	fmt.Println(sub)
	if err!=nil{
		panic(err)
	}

	stream, err := sqlite.Twitch_user_online(sub)

	if err!=nil{
		return err
	}


	for stream {
		//We have to check to see if there is an active prediction here that was set by me. IF it was not set by me, then we need to wait.
	prediction_id, _, err := sqlite.Get_predictions(user.User_id, "ACTIVE")

	//The problem is that I check my db to see if there are any active predictions. Right now there is not.
	//However, when I check tw

	if err!=nil{
		return err
	}
		own_active_prediction, err := twitch.Check_prediction(sub, user.Access_token, prediction_id)
		fmt.Println(own_active_prediction)
		if err!= nil{
			return err
		}

		 if own_active_prediction == "our_prediction" {
			fmt.Println("Our prediction")
			err = Watch_prediction(sub, user)

			if err!=nil{
				return err
			}

		 }else if own_active_prediction == "no_active_prediction"{
			fmt.Println("prediction started")
			err = twitch.Start_prediction(user)

			if err!=nil{
				panic(err)
			}
		} else if own_active_prediction == "not_our_prediction" {
			fmt.Println("Not our prediction")
			time.Sleep(30 * time.Second)
		} 

		if err!=nil{
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
func Watch_prediction(sub string, user sqlite.Twitch_user)error{

	prediction_id, created_at ,err := sqlite.Get_predictions(sub, "ACTIVE")

	if err!=nil{
		return err
	}

	if prediction_id == ""{
		return err
	}

	t_created_at, err := time.Parse(time.RFC3339, created_at)

	if err!=nil{
		panic(err)
	}

	new_battle := "no_new_battles"

	for (new_battle == "no_new_battles"){

		user, err := sqlite.Get_twitch_user("sub", sub)

		if err!=nil{return err}

		if user.Online != 1{
			return nil
		}

		new_battle, err := clash.New_battle(user.Player_tag, t_created_at)

		if err!=nil{
			return err
		}

		if new_battle == "no_new_battles"{
			fmt.Println("no new battles")
			time.Sleep(30 * time.Second)
		} else if new_battle == "win"{
			fmt.Println("Battle win")	
			outcome_id, err := sqlite.Get_prediction_outcome_id(prediction_id, 1)
			if err!=nil{
				return err
			}
			if outcome_id == ""{
				fmt.Println("Outcome id was blanck")
			}
			err = twitch.End_prediction(prediction_id, outcome_id, user.User_id, user.Access_token)
			if err !=nil{
				return err
			}
			return nil
		} else if new_battle == "lose"{
			outcome_id, err := sqlite.Get_prediction_outcome_id(prediction_id, 0)
			fmt.Println("Battle lose")	
			if err!=nil{
				return err
			}
			if outcome_id == ""{
				fmt.Println("Outcome id was blanck")
			}
			err = twitch.End_prediction(prediction_id, outcome_id, user.User_id, user.Access_token)
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