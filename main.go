package main

import (
	//"fmt"

	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	"fmt"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	//twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

func main(){
	//PUT CLIENT SECRET IN as second argument
	 //twitch.Get_user_info("oxalate", "")
	 //twitch.Generate_state()
	// twitch.Scope_requests("prediction")
	//twitch.Generate_authorize_app_url("", "prediction")
	//sqlite.Create_twitch_database()
	err := sqlite.Write_state("Testing State")

	if err != nil {
		fmt.Println(err)
		panic("Write Panic")
	}

	found, err := sqlite.Check_state("Testing State")
	if err!= nil{
		fmt.Println(err)
		panic("Check Panic")
	}
	fmt.Print(found)
}