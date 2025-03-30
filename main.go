package main

import (
	"errors"
	"fmt"

	//app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/server"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/ncruces/go-sqlite3"
	//twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	//logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
)

func main() {

	db, err := sqlite.Open_db("file:royale_bets") 
	defer db.Close()

	if err != nil {
		err = errors.New("FILE sqlite_helper FUNC: open_db CALL: sqlite3.Open " + err.Error())
		panic(err)
	}

	if err = sqlite.Create_twitch_database(db); err!=nil {
		panic(err)
	}

	if err = sqlite.Create_royale_bets_db(db); err != nil{
		panic(err)
	}


	// Env_struct, err := logger.Get_env_variables("/Users/jonathanlewis/Documents/Projects/clash_royale_twitch_prediction_bot/test.env")
	// if err != nil {
	// 	panic(err)
	// }
	//logger := logger.NewStandardLogger()
	//Write_test_user(Env_struct)
	test_clash_db(db)

	//server.Start_server(logger, Env_struct, db)
}

func test_db() {
	// err := sqlite.Create_twitch_database()
	// if err != nil{
	// 	panic(err)
	// }

	// user, err := sqlite.Get_twitch_user("sub","651008027")

	// if err!=nil{
	// 	panic(err)
	// }
	// fmt.Println(user)

	// status, err := twitch.Validate_token(user.Access_token, user.User_id)

	// if err!=nil{
	// 	fmt.Println(status)
	// }

}

func test_clash_db(db *sqlite3.Conn){

	viewers := []sqlite.Royale_bets_viewer{
		{1, "Streamer1", "ViewerA", 100, 1000000},
		{2, "Streamer2", "ViewerB", 200, 1000000},
	}
	for _, viewer := range viewers {
		if err := sqlite.Insert_royale_bets_viewer(db, viewer); err != nil {
			fmt.Println("Error inserting viewer:", err)
		}
	}

	streamers := []sqlite.Royale_bets_streamer{
		{1625097600, "Streamer1", 1625101200, 5, 2},
		{1625097700, "Streamer2", 1625101300, 3, 4},
	}
	for _, streamer := range streamers {
		if err := sqlite.Insert_royale_bets_streamer(db, streamer); err != nil {
			fmt.Println("Error inserting streamer:", err)
		}
	}

	results := []sqlite.Battle_result{
		{"Player1", 1625097600, 3, 1},
		{"Player2", 1625097700, 2, 2},
	}
	for _, result := range results {
		if err := sqlite.Insert_battle_result(db, result); err != nil {
			fmt.Println("Error inserting battle result:", err)
		}
	}

	results_print, err := sqlite.Get_battle_result(db, "Player1", 1225097700)
	results_print1, err := sqlite.Get_battle_result(db, "Player2", 1225097700)
	if err!= nil{panic(err)}
	fmt.Println(results_print)
	fmt.Println(results_print1)
}

func test_get_all_access_tokens(db *sqlite3.Conn) {
	Access_tokens, err := sqlite.Get_all_access_tokens(db)

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(Access_tokens); i++ {
		fmt.Println(Access_tokens[i])
	}
}

func test_twitch_api() {
	// url, err := twitch.Generate_authorize_app_url(App_id, "prediction")

	// if err != nil{
	// 	panic(err)
	// }
	// fmt.Println(url)

	//twitch.Start_prediction()

}
