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
			err = errors.New("there was a problem writing to the string scope maker")
			return "", err
		}

		if i+1 != len(scope_array){
			return_scope_string.WriteString(" ")
		}
	}
	return return_scope_string.String(), nil
}

// func Scope_packer() (string, error){
// 	return "", err
// }
