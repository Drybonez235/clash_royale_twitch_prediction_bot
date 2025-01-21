package main

import (
	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	"fmt"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

func main(){
	//PUT CLIENT SECRET IN as second argument
	 //twitch.Get_user_info("oxalate", "")
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
	//test_db()
	//test_twitch_api()
	test_twitch_api()
	twitch.Start_server()
}

func test_db(){
	// err := sqlite.Create_twitch_database()
	// if err != nil{
	// 	panic(err)
	// }

	err := sqlite.Write_state_nonce("Testing 1", "state")

	if err != nil{
		panic(err)
	}
	err = sqlite.Write_state_nonce("Testing nonce", "nonce")
	
	if err != nil{
		panic(err)
	}

	here, err := sqlite.Check_state_nonce("Invalid", "state")

	if err != nil{
		panic(err)
	}

	fmt.Println(here)

	state_here, err := sqlite.Check_state_nonce("Testing 1", "state")

	if err != nil{
		panic(err)
	}

	fmt.Println(state_here)

	nonce_here, err := sqlite.Check_state_nonce("Testing nonce", "nonce")

	if err != nil{
		panic(err)
	}

	fmt.Println(nonce_here)

}

func test_twitch_api(){
	url, err := twitch.Generate_authorize_app_url("now6dwkymg4vo236ius5d0sn82v9ul", "prediction")

	if err != nil{
		panic(err)
	}

	fmt.Println(url)

}