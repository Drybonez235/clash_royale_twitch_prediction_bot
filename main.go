package main

import (
	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
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
//const App_id ="now6dwkymg4vo236ius5d0sn82v9ul"
//const Secret = ""

func main(){
	
	
	test_db()

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

	// status, err := twitch.Validate_token("", user.User_id)

	// if err!=nil{
	// 	fmt.Println(status)
	// }

	// if !status{
	// 	refreshed, err := twitch.Refresh_token(user.Refresh_token, user.User_id)

	// 	if !refreshed || err !=nil{
	// 		panic(err)
	// 	}
	// 	user, err = sqlite.Get_twitch_user("sub", user.User_id)
	// 	if !refreshed || err !=nil{
	// 		panic(err)
	// 	}
	// }
	
	status, err := twitch.Validate_token(user.Access_token, user.User_id)

	if err!=nil{
		fmt.Println(status)
	}

	// fmt.Println(status)

	// err = twitch.Start_prediction(user.Access_token, user.User_id, user.Display_Name)
	// if err!=nil{
	// 	panic(err)
	// }	


}

func test_twitch_api(){
	url, err := twitch.Generate_authorize_app_url(App_id, "prediction")

	if err != nil{
		panic(err)
	}
	fmt.Println(url)


}