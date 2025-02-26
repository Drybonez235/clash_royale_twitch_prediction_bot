package server

import "github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"

func Validate_all_tokens()(error){
	var Token_list []string

	Token_list, err := sqlite.Get_all_tokens()

	for i:=0; i<len(Token_list); i++{
		token := Token_list[i]	
		err := Validate_one_token(token)
		if err!=nil{
			return err
		}
	}
	return err
}

func Validate_one_token(token string)(error){
	if
}