package server

import ()


//These are for receiving subscription alerts from Twitch.
type Challenge_struct struct {
	Challenge   string      `json:"challenge"`
	Subscription Subscription `json:"subscription"`
}

type Subscription struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	Type      string     `json:"type"`
	Version   string     `json:"version"`
	Cost      int        `json:"cost"`
	Condition Condition  `json:"condition"`
	Transport Transport  `json:"transport"`
	CreatedAt string     `json:"created_at"`
}

type Condition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
}



//This is for webhook body requests
type WebhookEvent struct {
	Type      string     `json:"type"`
	Version   string     `json:"version"`
	Condition Condition_User_ID  `json:"condition"`
	Transport Transport_Secret  `json:"transport"`
}

type Condition_User_ID struct {
	UserID string `json:"user_id"`
}

type Transport_Secret struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}