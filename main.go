package main

import (
	"fmt"
	//app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/server"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	//twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

//So at this point I am trying to set up a system that checks to make siure the token is valid. IF it isn't valid, I need to request a new token using a refresh token. Then I am wrtiing
//To the db the new access token and refresh token. Something somewhere is breaking.

const App_id ="b2109dc3a41733acaa7b3fa355df4c" //Test app id
const Secret = "dacb3721ea3023f1e955a053d91f24" //Test secret
const user_id = "29277192"

func main(){
	// var1, err := twitch.Get_display_name("29277192")

	// if err!=nil{panic(err)}
	// fmt.Println(var1)

	var2, err := clash.String_time_to_time_time("20240101T0202002Z")
	if err!=nil{panic(err)}
	fmt.Println(var2)
	test_test_twitch_api()

	test_get_all_access_tokens()

	//test_db()
 	//server.Start_server()
	//test_event_sub()
	//test_app()
	//test_test_twitch_api()
}

func test_db(){
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

func test_get_all_access_tokens(){
	Access_tokens, err := sqlite.Get_all_access_tokens()

	if err!=nil{
		panic(err)
	}

	for i:=0; i < len(Access_tokens); i++{
		fmt.Println(Access_tokens[i])
	}
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

	err := sqlite.Write_twitch_info("29277192", "Name", "b9200d33b7b6611", "ghsifnwofieflakdjenfonf", "not important", "bearer", "","",0,0,"",0,"2VL9VP8Y0")

	if err!=nil{
		panic(err)
	}
	err = sqlite.Write_twitch_info("10209020", "Name", "hgjofhqofn", "ig0pqidhduchauhd", "not important", "bearer", "","",0,0,"",0,"2VL9VP8Y0")

	if err!=nil{
		panic(err)
	}

	// user, err := sqlite.Get_twitch_user("sub", user_id)	

	// if err!= nil{
	// 	fmt.Println(err)
	// }
	// err = server.Create_EventSub(user, "stream.online")
	// if err!=nil{
	// 	panic(err)
	// }

}


// func test_app(){
// 	user, err := sqlite.Get_twitch_user("sub", user_id)	

// 	if err!= nil{
// 		fmt.Println(err)
// 	}

// 	err = app.Start_prediction_app(user.User_id)

// 	if err!=nil{
// 		panic(err)
// 	}
// }

func test_event_sub(){
	user, err := sqlite.Get_twitch_user("sub", user_id)	

	if err!= nil{
		fmt.Println(err)
	}

	err = server.Create_EventSub(user, "stream.online")

	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
}