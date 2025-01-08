package twitch_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func request_oath_token(client_id string, client_secret string, streamer_name string) (string, error) {
	fmt.Println("This is the request auth token function")
	twitch_oath_url := "https://id.twitch.tv/oauth2/token?" //{"https://id.twitch.tv/oauth2/token"}
	//urlencoded_string := "client_id="+client_id+"&client_secret="+client_secret+"&grant_type=client_credentials"
	//url_map := url.Values{}
	url_quary := url.Values{}
		
	url_quary.Add("Client_id", client_id)
	url_quary.Add("Client_secret", client_secret)
	url_quary.Add("grant_type", "client_credentials")
	url_encoded_string := url_quary.Encode()
	
	fmt.Println(twitch_oath_url+url_encoded_string)
	//fmt.Println(url_encoded_string)



	
	client := &http.Client{}
	
	req, err := http.NewRequest("Post", twitch_oath_url + url_encoded_string, nil)

	if err != nil {
		fmt.Print(err)
		return "Hmm", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	  }
	
	if err != nil || resp.Status != "200 OK"{
		fmt.Print(err)
		return "Hmm", err
	}

	fmt.Println("This is after we defer and close the respone body")

	body, err := io.ReadAll(resp.Body)

	fmt.Println(body)

	if err != nil{
		fmt.Print(err)
		return "Hmm", err
	}

	access_token_json, err := data_handler_array(body)

	if err != nil {
		fmt.Print(err)
		return "Hmm", err	
	}
	var access_token_string = string(access_token_json[0])

	call_user_endpoint(streamer_name, access_token_string)
	return resp.Status, err
}

func data_handler_array(data []byte) ([3]string, error){
	fmt.Println("This is the Data handler function")
	var json_array []map[string]any
	var access_token_array [3]string

	err := json.Unmarshal(data, &json_array)
	
	if(err != nil){
		err = errors.New("something went wrong")
	}
	zero := json_array[0] 
	
	access_token_array[0] =  zero["access_token"].(string)
	access_token_array[1] =  zero["expires_in"].(string)	
	access_token_array[2] =  zero["token_type"].(string)

	fmt.Println("The access token is" + access_token_array[0])
 
	return access_token_array, err
}

func call_user_endpoint(streamer_id string, access_token string) (string, error) {
	fmt.Println("This is the call user endpoint function")
	twitch_get_user_endpoint := "https://api.twitch.tv/helix/users?login=twitchdev"
	bearer_token := "Bearer " + access_token

	client := &http.Client{}
	
	req, err := http.NewRequest("Post", twitch_get_user_endpoint, nil)

	req.Header.Add("Authorization", bearer_token)
	req.Header.Add("Client-Id", "wbmytr93xzw8zbg0p1izqyzzc5mbiz")

	resp, err := client.Do(req)
	
	if err != nil || resp.Status != "200 OK"{
		fmt.Print(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println(body) 

	return "Do we get here?", err
}

func Get_user_info(streamer_id string){
	client_id := "now6dwkymg4vo236ius5d0sn82v9ul"
	client_secret := "j4hoyew3efaptladimr7xgtzu1802d"
	request_oath_token(client_id, client_secret, streamer_id)
} 
