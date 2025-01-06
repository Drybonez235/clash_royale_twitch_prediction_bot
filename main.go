package main

import (
	"fmt"

	clash "github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api"
	twitch "github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)

func main(){
	message1 := twitch.Get_twitch()
	message2 := clash.Get_clash()
	fmt.Println(message1)
	fmt.Println(message2)
}