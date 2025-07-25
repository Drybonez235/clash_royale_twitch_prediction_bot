package main

import (
	"errors"
	"fmt"

	//app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/server"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/ncruces/go-sqlite3"
	//twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
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

	//test_clash_db(db)

	Env_struct, err := logger.Get_env_variables("./.env")
	if err != nil {
		panic(err)
	}
	logger := logger.NewStandardLogger()
	
	sqlite.Get_all_battle_results(db)
	server.Start_server(logger, Env_struct, db)
}

func test_db() {
	

}

func test_clash_db(db *sqlite3.Conn){

	// viewers := []sqlite.Royale_bets_viewer{
	// 	{10, "Streamer1", "ViewerA", 100, 1},
	// 	{11, "Streamer2", "ViewerB", 200, 1},
	// }
	// for _, viewer := range viewers {
	// 	if err := sqlite.Insert_royale_bets_viewer(db, viewer); err != nil {
	// 		fmt.Println("Error inserting viewer:", err)
	// 	}
	// }

// 	streamers := []sqlite.Royale_bets_streamer{
// 		{0, "2VL9VP8Y0", 0, 5, 2},
// 	}
// 	for _, streamer := range streamers {
// 		if err := sqlite.Insert_royale_bets_streamer(db, streamer); err != nil {
// 			fmt.Println("Error inserting streamer:", err)
// 		}
// 	}

// 	results := []sqlite.Battle_result{
// 		{"2VL9VP8Y0", 0, 3, 1},
// 		{"2VL9VP8Y0", 1, 2, 2},
// 	}
// 	for _, result := range results {
// 		if err := sqlite.Insert_battle_result(db, result); err != nil {
// 			fmt.Println("Error inserting battle result:", err)
// 		}
// 	}
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
