package main

import (
	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	"fmt"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

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

const secret = ""

func main(){
	test_db()
	 //twitch.Generate_state()
	// twitch.Scope_requests("prediction")
	//url, err := twitch.Generate_authorize_app_url("now6dwkymg4vo236ius5d0sn82v9ul", "prediction")


	// if err != nil {
	// 	fmt.Println(err)
	// 	panic("problem with url generator")
	// }

	// found, err := sqlite.Check_state("Testing State")

	// if err!= nil{
	// 	fmt.Println(err)
	// 	panic("Check Panic")
	// }
	//err = twitch.Start_server()

	// if err != nil{
	// 	panic("there was a problem with the server")
	// }
	
	// test_twitch_api()
	// twitch.Start_server()

}

func test_db(){
	// err := sqlite.Create_twitch_database()
	// if err != nil{
	// 	panic(err)
	// }

	// var tu Twitch_user_info 

	// tu.sub = "122222"

	// tu.access_token = "19385930282"
	// tu.display_name = "oxalate"
	// tu.refresh_token = "refresh token"
	// tu.scope = "scope"
	// tu.token_type = "bearer"
	// tu.app_request = "clash_royal_prediction_bot"
	// tu.app_received = "Clash royalk_prediction_bot"
	// tu.token_exp = 50000.0
	// tu.token_iat = 12345.0
	// tu.token_iss = "Twitch"

	// err := sqlite.Write_twitch_info(tu.sub, tu.display_name, tu.access_token, tu.refresh_token, tu.scope, tu.token_type, tu.app_request, tu.app_received, tu.token_exp, tu.token_iat, tu.token_iss)

	// if err != nil{
	// 	panic(err)
	// }

	err := sqlite.Get_twitch_user("sub","651008027")
	if err!=nil{
		panic(err)
	}
	// err := sqlite.Write_state_nonce("Testing 1", "state")

	// if err != nil{
	// 	panic(err)
	// }
	// err = sqlite.Write_state_nonce("Testing nonce", "nonce")
	
	// if err != nil{
	// 	panic(err)
	// }

	// here, err := sqlite.Check_state_nonce("Invalid", "state")

	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Println(here)

	// state_here, err := sqlite.Check_state_nonce("Testing 1", "state")

	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Println(state_here)

	// nonce_here, err := sqlite.Check_state_nonce("Testing nonce", "nonce")

	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Println(nonce_here)

}

func test_twitch_api(){
	url, err := twitch.Generate_authorize_app_url("now6dwkymg4vo236ius5d0sn82v9ul", "prediction")

	if err != nil{
		panic(err)
	}

	fmt.Println(url)
	//display, err := twitch.Get_display_name("Oxalate")

	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Print(display)

}