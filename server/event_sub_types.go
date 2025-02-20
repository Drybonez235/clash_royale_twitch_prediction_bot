package server

import ()

//When an event that you subscribe to occurs, Twitch sends your event handler 
// a notification message that contains the event’s data.
type WebhookNotification struct {
	Subscription Subscription `json:"subscription"`
	Event        Event        `json:"event"`
}

type Subscription struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	Type      string     `json:"type"`
	Version   string     `json:"version"`
	Cost      int        `json:"cost"`
	Condition Condition  `json:"condition"`
	Transport Transport  `json:"transport"`
	CreatedAt string  `json:"created_at"`
}

type Condition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
}

type Event struct {
	UserID               string    `json:"user_id"`
	UserLogin            string    `json:"user_login"`
	UserName             string    `json:"user_name"`
	BroadcasterUserID    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	FollowedAt           string 	`json:"followed_at"`
}
//End of Subscription event type

//When you successfully subscribe to an event and respond to the challenge. 
//Creates an EventSub subscription.
type EventSubResponse struct{
	Data []Subscription
	Total        int           `json:"total"`
	TotalCost    int           `json:"total_cost"`
	MaxTotalCost int           `json:"max_total_cost"`
}


//Challenge request type.
type Challenge_struct struct {
	Challenge string `json:"challenge"`
	Subscription Subscription `json:"subscription"`
}

//This is the json type for the subrequest json body. Not the response or anything.
type EventSubRequest struct {
	Type      string     `json:"type"`
	Version   string     `json:"version"`
	Condition ConditionSubRequest  `json:"condition"`
	Transport TransportSubRequest  `json:"transport"`
}

type ConditionSubRequest struct {
	UserID string `json:"user_id"`
}

type TransportSubRequest struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}



//These are for receiving subscription alerts from Twitch.
// type Challenge_struct struct {
// 	Type      string     `json:"type"`
// 	Version   string     `json:"version"`
// 	Condition Condition  `json:"condition"`
// 	Transport Transport_Secret  `json:"transport"`
// }

// type Subscription struct {
// 	ID        string     `json:"id"`
// 	Status    string     `json:"status"`
// 	Type      string     `json:"type"`
// 	Version   string     `json:"version"`
// 	Cost      int        `json:"cost"`
// 	Condition Condition  `json:"condition"`
// 	Transport Transport  `json:"transport"`
// 	CreatedAt string     `json:"created_at"`
// }

// type Condition struct {
// 	BroadcasterUserID string `json:"broadcaster_user_id"`
// }

// type Transport struct {
// 	Method   string `json:"method"`
// 	Callback string `json:"callback"`
// }

// //Receive webhook body json
// type Transport_Secret struct {
// 	Method   string `json:"method"`
// 	Callback string `json:"callback"`
// 	Secret   string `json:"secret"`
// }

// type Condition_User_ID struct {
// 	UserID string `json:"user_id"`
// }

// type Data struct {
// 	ID         string    `json:"id"`
// 	Status     string    `json:"status"`
// 	Type       string    `json:"type"`
// 	Version    string    `json:"version"`
// 	Condition  Condition `json:"condition"`
// 	CreatedAt  string `json:"created_at"`
// 	Transport  Transport `json:"transport"`
// 	Cost       int       `json:"cost"`
// }

// type Webhook_body struct {
// 	Data        []Data `json:"data"`
// 	Total       int    `json:"total"`
// 	TotalCost   int    `json:"total_cost"`
// 	MaxTotalCost int   `json:"max_total_cost"`
// }