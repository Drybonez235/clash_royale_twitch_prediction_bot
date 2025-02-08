package main

import (
	clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	"fmt"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

//So at this point I am trying to set up a system that checks to make siure the token is valid. IF it isn't valid, I need to request a new token using a refresh token. Then I am wrtiing
//To the db the new access token and refresh token. Something somewhere is breaking.

type Twitch_user_info struct{
	sub string
	display_name string
	access_token string
	refresh_token string
	scope string
	token_type string
	app_request string
	app_received string
	token_exp float64
	token_iat float64
	token_iss string
}
const App_id ="b2109dc3a41733acaa7b3fa355df4c" //Test app id
const Secret = "dacb3721ea3023f1e955a053d91f24" //Test secret
const user_id = "29277192"

func main(){
	//Test_clash_api()
	test_test_twitch_api()
	// user, err := sqlite.Get_twitch_user("sub", user_id)	

	// if err!= nil{
	// 	panic(err)
	// }
	
	// err = twitch.Start_prediction(user)

	// if err!= nil{
	// 	panic(err)
	// }
	// test_twitch_api()
	// twitch.Start_server()

}

func test_db(){
	// err := sqlite.Create_twitch_database()
	// if err != nil{
	// 	panic(err)
	// }

	user, err := sqlite.Get_twitch_user("sub","651008027")

	if err!=nil{
		panic(err)
	}
	fmt.Println(user)

	// status, err := twitch.Validate_token(user.Access_token, user.User_id)

	// if err!=nil{
	// 	fmt.Println(status)
	// }

	// if !status{
	// refreshed, err := twitch.Refresh_token(user.Refresh_token, user.User_id)

	// 	if !refreshed || err !=nil{
	// 		panic(err)
	// 	}
	// fmt.Println(refreshed)
	// 	user, err = sqlite.Get_twitch_user("sub", user.User_id)
	// 	if !refreshed || err !=nil{
	// 		panic(err)
	// 	}
	// }

	// fmt.Println(status)

	// err = twitch.Start_prediction(user.Access_token, user.User_id, user.Display_Name)
	// if err!=nil{
	// 	panic(err)
	// }	
	// sqlite.Remove_twitch_user(user_id)

}

func test_twitch_api(){
	// url, err := twitch.Generate_authorize_app_url(App_id, "prediction")

	// if err != nil{
	// 	panic(err)
	// }
	// fmt.Println(url)

	//twitch.Start_prediction()

}

func test_test_twitch_api(){
	// err := twitch.Test_request_user_oath_token(user_id)
	// if err != nil{
	// 	panic(err)
	// }

	// err := sqlite.Write_twitch_info(user_id, "Name", "13e35cf9690bc85", "", "not important", "bear", "","",0,0,"")

	// if err!=nil{
	// 	panic(err)
	// }

	user, err := sqlite.Get_twitch_user("sub", user_id)	

	if err!= nil{
		fmt.Println(err)
	}

	err = twitch.Start_prediction(user)

	if err!= nil{
		panic(err)
	}

	prediction_id, err := sqlite.Get_predictions(user.User_id, "ACTIVE")

	if err !=nil{
		panic(err)
	}

	outcome, err := sqlite.Get_prediction_outcome_id(prediction_id, 1)

	if err !=nil{
		panic(err)
	}

	fmt.Println(outcome)

	err = twitch.End_prediction(prediction_id, outcome, user.User_id, user.Access_token)

	if err!=nil{
		panic(err)
	}

}

func Test_clash_api(){
	player_tag := "2YJRUQ2Q"
	battles, err := clash.Get_prior_battles(player_tag)

	if err!=nil{
		panic(err)
	}
	// fmt.Println(battles.Matches[0].Team)

	battle_0 := battles.Matches[0].BattleTime
	
	time, err := clash.String_time_to_time_time(battle_0)

	fmt.Print(time)
}