package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
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
			stream_start(webhook_struct.Event.BroadcasterUserID)
		}
	}
	
	return nil
}

func stream_start(streamer_id string)(error){
	fmt.Println("Stream started fired")
	user, err := sqlite.Get_twitch_user("sub", streamer_id)

	if err!=nil{return err}

	err = app.Start_prediction_app(user.User_id)

	if err!=nil{
		return err
	}

	return nil
}