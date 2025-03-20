package twitch_api

import (
	"strings"
	"errors"
)

func Scope_unpacker(scope_array []string) (string, error){
	var return_scope_string strings.Builder

	for i := 0; i < len(scope_array); i++{
		scope_string := scope_array[i]
		
		_, err := return_scope_string.WriteString(scope_string)
		if err!= nil{
			err = errors.New("FILE: Twitch_scope_handler FUNC: Scope_unpacker CALL: string.WriteString" + err.Error())
			return "", err
		}

		if i+1 != len(scope_array){
			return_scope_string.WriteString(" ")
		}
	}
	return return_scope_string.String(), nil
}

func Scope_requests(request string) (string, error){
	var err error 

	request_string := ""

	if request == "prediction"{
		request_string = "channel:manage:predictions openid" 
		return request_string, err
	} else {
		err = errors.New("invalid scope request")
	}

	if err != nil {
		err = errors.New("FILE: Twitch_scope_handler FUNC: Scope_requests " + err.Error())
		return "", err
	}
	return request, err
}

