package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)
	
func Handle_event(w http.ResponseWriter, req *http.Request)(error){
	fmt.Println("We are hadnling the event!")
	var webhook_struct WebhookNotification

	req_body, err := io.ReadAll(req.Body)

	if err!=nil{
		return err
	}

	err = json.Unmarshal(req_body, &webhook_struct)

	if err!=nil{
		return err
	}
	exists, err := sqlite.Get_sub_event(webhook_struct.Subscription.ID)

	if err!=nil{
		return err
	}
	if exists {
		return nil
	} else {
		err = sqlite.Write_sub_event(webhook_struct.Subscription.ID)
		if err!=nil{
			return err
		}
		//Logic goes here.
		if webhook_struct.Subscription.Type == "stream.online"{
			fmt.Println(webhook_struct.Event.BroadcasterUserID)
			stream_start(webhook_struct.Event.BroadcasterUserID)
		} else if webhook_struct.Subscription.Type == "stream.offline"{
			stream_end(webhook_struct.Event.BroadcasterUserID)
		}
	}
	
	return nil
}

func stream_start(streamer_id string)(error){

	if streamer_id == ""{
		return errors.New("there was no streamer id")
	}
	fmt.Println("Stream started fired")
	fmt.Println(streamer_id)
	user, err := sqlite.Get_twitch_user("sub", streamer_id)


	if err!=nil{return err}

	if user.User_id == ""{
		err = errors.New("there was no streamer associated with the sub id pulled from the request")
		return err
	}

	err = sqlite.Update_online(streamer_id, 1)

	if err!=nil{
		return err
	}

	fmt.Println("This is the user id")
	fmt.Println(user.User_id)
	go app.Start_prediction_app(user.User_id)

	if err!=nil{
		return err
	}

	return nil
}

//Not tested yet
func stream_end(streamer_id string)(error){
	var err error

	if streamer_id == ""{
		err = errors.New("streamer_id was blank")
		return err
	}

	err = sqlite.Update_online(streamer_id, 0)

	if err!=nil{return err}

	user, err := sqlite.Get_twitch_user("sub", streamer_id)

	if err!=nil{return err}

	prediction_id, _ ,err := sqlite.Get_predictions(streamer_id, "ACTIVE")

	if err!=nil{return err}	

	err = twitch_api.Cancel_prediction(prediction_id, streamer_id, user.Access_token)

	if err!=nil{return err}

	return nil
}