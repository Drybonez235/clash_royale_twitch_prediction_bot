package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	app "github.com/Drybonez235/clash_royale_twitch_prediction_bot/app"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/logger"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/sqlite"
	"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
)
	
func Handle_event(w http.ResponseWriter, req *http.Request, logger *logger.StandardLogger)(error){
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

	//Exists deals with duplicated events
	if exists {
		return nil
	} else {
		err = sqlite.Write_sub_event(webhook_struct.Subscription.ID)
		if err!=nil{
			return err
		}
		
		if webhook_struct.Subscription.Type == "stream.online"{
			logger.Info("stream online for: " + webhook_struct.Event.BroadcasterUserLogin)
			err = stream_start(webhook_struct.Event.BroadcasterUserID)
			if err!=nil{
				return err
			}
		} else if webhook_struct.Subscription.Type == "stream.offline"{
			logger.Info("stream offline for: " + webhook_struct.Event.BroadcasterUserLogin)
			err = stream_end(webhook_struct.Event.BroadcasterUserID)
			if err!=nil{
				return err
			}
		}
	}
	return nil
}

func stream_start(streamer_id string)(error){

	if streamer_id == ""{
		return errors.New("stream start: streamer id was blank")
	}
	user, err := sqlite.Get_twitch_user("sub", streamer_id)

	if err!=nil{return err}

	if user.User_id == ""{
		err = errors.New("stream start: streamer not found in db")
		return err
	}

	err = sqlite.Update_online(streamer_id, 1)

	if err!=nil{
		return err
	}

	go app.Start_prediction_app(user.User_id)

	return nil
}

func stream_end(streamer_id string)(error){
	var err error

	if streamer_id == ""{
		err = errors.New("stream end: streamer id was blank")
		return err
	}

	err = sqlite.Update_online(streamer_id, 0)

	if err!=nil{return err}

	user, err := sqlite.Get_twitch_user("sub", streamer_id)

	if err!=nil{return err}

	if user.User_id == ""{
		err = errors.New("stream end: streamer not found in db")
		return err
	}


	prediction_id, _ ,err := sqlite.Get_predictions(streamer_id, "ACTIVE")

	if err!=nil{
		return err
	}

	if prediction_id == ""{
		return nil
	}

	err = twitch_api.Cancel_prediction(prediction_id, streamer_id, user.Access_token)

	if err!=nil{return err}

	return nil
}