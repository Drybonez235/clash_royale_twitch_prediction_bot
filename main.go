package main

import (
	"fmt"
	//app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	//clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/server"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	//twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	logger "github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
)

func main() {

	Env_struct, err := logger.Get_env_variables("/Users/jonathanlewis/Documents/Projects/clash_royale_twitch_prediction_bot/.env")
	if err != nil {
		panic(err)
	}
	logger := logger.NewStandardLogger()
	server.Start_server(logger, Env_struct)
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

func test_get_all_access_tokens() {
	Access_tokens, err := sqlite.Get_all_access_tokens()

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

func test_test_twitch_api() {
	// err := twitch.Test_request_user_oath_token(user_id)
	// if err != nil{
	// 	panic(err)
	// }

	err := sqlite.Write_twitch_info("29277192", "Name", "9d4bc3bbde86b87", "ghsifnwofieflakdjenfonf", "not important", "bearer", "", "", 0, 0, "", 0, "2VL9VP8Y0")

	if err != nil {
		panic(err)
	}
	// err = sqlite.Write_twitch_info("10209020", "Name", "hgjofhqofn", "ig0pqidhduchauhd", "not important", "bearer", "","",0,0,"",0,"2VL9VP8Y0")

	// if err!=nil{
	// 	panic(err)
	// }

	// user, err := sqlite.Get_twitch_user("sub", user_id)

	// if err!= nil{
	// 	fmt.Println(err)
	// }
	// err = server.Create_EventSub(user, "stream.online")
	// if err!=nil{
	// 	panic(err)
	// }

}
