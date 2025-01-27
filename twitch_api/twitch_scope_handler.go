package twitch_api

import (
	"strings"
	"errors"
	"fmt"
)

func Scope_unpacker(scope_array []string) (string, error){
	var return_scope_string strings.Builder

	for i := 0; i < len(scope_array); i++{
		scope_string := scope_array[i]
		
		_, err := return_scope_string.WriteString(scope_string)
		if err!= nil{
			err = errors.New("there was a problem writing to the string scope maker")
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
		return "", err
	}

	fmt.Println(request)
	return request, err
}
// func Scope_packer() (string, error){
// 	return "", err
// }
