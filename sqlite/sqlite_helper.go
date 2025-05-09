package sqlite

import (
	"errors"
	"fmt"
	//"log"

	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const file = "file:twitch_authorization"

type Twitch_user struct{
	User_id string
	Display_Name string
	Access_token string
	Refresh_token string
	Scope string
	Token_type string
	App_request string
	App_received string
	Token_exp int
	Token_iat int
	Token_iss string
	Online int
	Player_tag string
}

type Twitch_user_refresh struct{
	User_id string
	Access_token string
	Refresh_token string
}

func Open_db(file string)(*sqlite3.Conn, error){
	db, err := sqlite3.Open(file) 
	if err != nil {
		err = errors.New("FILE sqlite_helper FUNC: open_db CALL: sqlite3.Open " + err.Error())
		return db, err
	}	
	return db, nil
}

func Create_twitch_database(db *sqlite3.Conn) error {

	err := db.Exec(`CREATE TABLE IF NOT EXISTS state (state_value text, player_id text)`)
	if err != nil{
		err = errors.New("db: there was a problem creating the state table")
		return err
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS twitch_user_info (sub text, display_name text, access_token text, refresh_token text, scope text, token_type text, app_request text,
	app_received text, token_exp float, token_iat float, token_iss text, online int, clash_royale_player_tag text)`)

	if err!=nil{
		return err
	}
	err = db.Exec(`CREATE TABLE IF NOT EXISTS prediction (broadcaster_id text, prediction_id text, status text, created_at text)`)

	if err != nil{
		err = errors.New("db: there was a problem creating the twitch_user_info table")
		return err
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS outcomes (prediction_id text, outcome_id text, title text, lose_win int)`)

	if err!=nil{
		err = errors.New("db: there was a problem creating the outcomes table")
		return err
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS Sub_Events (Sub_Event_ID text)`)

	if err!=nil{
		err = errors.New("db: there was a problem creating the sub_events table")
		return err
	}
	return err
}

func Write_state_nonce(db *sqlite3.Conn, state_nonce string, table string) error {
	
	header := ""

	if table == "state"{
		header = "state_value"
	} else if table == "nonce" {
		header = "nonce_value"
	} else {
		err := errors.New("FILE: sqlite_helper FUNC: Write_state_nonce INVALID PARAMETER: Must be state or nonce")
		return err
	}

	sql_command := fmt.Sprintf(`INSERT INTO '%s' ('%s') VALUES ('%s')`, table, header, state_nonce)
	err := db.Exec(sql_command) //`INSERT INTO state (state_value) VALUES ('Testing')`
	if err != nil{ 
		err = errors.New("FILE: sqlite_helper FUNC: Write_state_nonce CALL: db.Exec " + err.Error())
		return err
	}

	return nil
}

//This returns the state or nonce that was used to create the secure url and the player tag associated with that session.
func Check_state_nonce(db *sqlite3.Conn, state_nonce string, table string) (bool, string, error){
	
	header := ""

	if table == "state"{
		header = "state_value"
	} else if table == "nonce" {
		header = "nonce_value"
	} else {
		err := errors.New("FILE: sqlite_helper FUNC: Check_state_nonce INVALID PARAMETER: Must be state or nonce")
		return false, "", err
	}

	sql_query_string := fmt.Sprintf(`SELECT * FROM '%s' WHERE %s == '%s'`, table, header, state_nonce)
	sql_query, _, err := db.Prepare(sql_query_string)
	if err != nil{
		err = errors.New("FILE: sqlite_helper FUNC: Write_state_nonce CALL: db.Prepare " + err.Error())
		return false,"",err
	}
	if sql_query.Step() {
		player_tag := sql_query.ColumnText(1)
		sql_query.Close()
		err = delete_state_nonce(db, state_nonce, table)
		if err!=nil{return true, player_tag, err}
		return true, player_tag, nil
	} 
	
	return false,"",err
}

func delete_state_nonce(db *sqlite3.Conn, state_nonce string, table string) error {
	header := ""

	if table == "state"{
		header = "state_value"
	} else if table == "nonce" {
		header = "nonce_value"
	} else {
		err := errors.New("FILE: sqlite_helper FUNC: delete_state_nonce INVALID PARAMETER: Must be state or nonce")
		return err
	}

	sql_query_string := fmt.Sprintf(`DELETE FROM '%s' WHERE '%s' == '%s'`, table, header, state_nonce )
	err := db.Exec(sql_query_string)
	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: delete_state_nonce CALL: db.Exec " + err.Error())
		return err
	}
	
	return nil
}

//Works with the new fields!
func Write_twitch_info(db *sqlite3.Conn, sub string, display_name string, access_token string, refresh_token string, scope string, token_type string,
	 app_request string, app_received string, token_exp int, token_iat int, token_iss string, online int, player_tag string) error {
		err := Remove_twitch_user(db, sub)
		if err!=nil{
			return err
		}
		
		//(sub text, access_token text, refresh_token text, scope text, token_type text, app_request text
			//app_received text, token_exp float, token_iat float, token_iss text, online, clash_royale_player_tag)
		sql_table_values := "'sub', 'display_name', 'access_token', 'refresh_token', 'scope', 'token_type', 'app_request', 'app_received', 'token_exp', 'token_iat', 'token_iss', 'online', clash_royale_player_tag"
		sql_command := fmt.Sprintf("INSERT INTO twitch_user_info (%s) VALUES ('%s', '%s' , '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', %d, '%s')", sql_table_values, sub, display_name, access_token, refresh_token, scope, token_type, app_request, app_received, token_exp, token_iat, token_iss, online, player_tag)
		err = db.Exec(sql_command)	
		if err!=nil{
			err = errors.New("FILE: sqlite_helper FUNC: Write_twitch_info CALL: db.Exec " + err.Error())
			return err
		}
	return nil
}

func Get_twitch_user(db *sqlite3.Conn, id_type string, id string) (Twitch_user, error){
	var twitch_user Twitch_user

	field := ""
	if id_type == "sub"{
		field = "sub"
	} else if id_type == "display_name" {
		field = "display_name"
	} else {
		err := errors.New("FILE: sqlite_helper FUNC: Get_twitch_user INVALID PARAMETER: Must be sub or display_name")
		return twitch_user, err
	}
	sql_query_string := fmt.Sprintf(`SELECT * FROM twitch_user_info WHERE %s == '%s'`, field, id)
	sql_query, _, err := db.Prepare(sql_query_string)
	if err != nil{
		err = errors.New("FILE: sqlite_helper FUNC: Get_twitch_user CALL: db.Prepare " + err.Error())
		return twitch_user, err
	}

	for sql_query.Step() {
		twitch_user.User_id = sql_query.ColumnText(0)
		twitch_user.Display_Name= sql_query.ColumnText(1)
		twitch_user.Access_token= sql_query.ColumnText(2)
		twitch_user.Refresh_token= sql_query.ColumnText(3)
		twitch_user.Scope= sql_query.ColumnText(4)
		twitch_user.Token_type= sql_query.ColumnText(5)
		twitch_user.App_request= sql_query.ColumnText(6)
		twitch_user.App_received= sql_query.ColumnText(7)
		twitch_user.Token_exp= sql_query.ColumnInt(8)
		twitch_user.Token_iat= sql_query.ColumnInt(9)
		twitch_user.Token_iss= sql_query.ColumnText(10)
		twitch_user.Online = sql_query.ColumnInt(11)
		twitch_user.Player_tag = sql_query.ColumnText(12)
	}
	//PLEASE DO NOT FORGET ABOUT THIS...
	if twitch_user.User_id == ""{
		twitch_user.User_id = "29277192"
		return twitch_user, err
	}

	return twitch_user, nil	
}

func Twitch_user_online(db *sqlite3.Conn, sub string)(bool, error){

	sql_query_string := fmt.Sprintf(`SELECT online FROM twitch_user_info WHERE sub=='%s'`, sub)

	sql_query, _, err := db.Prepare(sql_query_string)	

	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Twitch_user_online CALL: db.Prepare " + err.Error())
		return false, err
	}
	false_true := 3
	for sql_query.Step(){
		false_true = sql_query.ColumnInt(0)
	}

	if false_true == 0{
		return false, nil
	} else if false_true == 1{
		return true, nil
	} else {
		err = errors.New("FILE: sqlite_helper FUNC: Twitch_user_online BUG: online not set for twitch user")
		return false, err
	}
}

func Update_online(db *sqlite3.Conn, sub string, online int)(error){

	if online >= 2 {
		err := errors.New("FILE: sqlite_helper FUNC: Update_online BUG: Online must be set to 0 or 1")
		return err
	}

	sql_query_string := fmt.Sprintf(`UPDATE twitch_user_info SET online = %d WHERE sub=='%s'`,online, sub)
	err := db.Exec(sql_query_string)

	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Update_online CALL: db.Exec " + err.Error())
		return err
	}
	return nil
}

func Update_tokens(db *sqlite3.Conn, access_token string, refresh_token string, sub string)(error){

	
	sql_query_string := fmt.Sprintf(`UPDATE twitch_user_info SET access_token = '%s', refresh_token = '%s' WHERE sub=='%s'`, access_token, refresh_token, sub)
	err := db.Exec(sql_query_string)
	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Update_tokens CALL: db.Exec " + err.Error())
		return err
	}
	return nil
}

func Remove_twitch_user(db *sqlite3.Conn, sub string) error {
	sql_string := fmt.Sprintf("DELETE FROM twitch_user_info WHERE sub=='%s'", sub)
	err := db.Exec(sql_string)
	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Remove_twitch_user CALL: db.Exec " + err.Error())
		return err
	}
	return nil
}

func Write_new_prediction(db *sqlite3.Conn, streamer_id string, prediction_id string, created_at string) error {
	
	sql_query_string := fmt.Sprintf(`INSERT INTO prediction ('broadcaster_id', 'prediction_id', 'status', 'created_at') VALUES('%s','%s','ACTIVE', '%s')`,streamer_id, prediction_id, created_at)
	err := db.Exec(sql_query_string)
	if err !=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Write_new_prediction CALL: db.Exec " + err.Error())
		return err
	}
	return nil
}

//Maybe instead we can make it accept the array of predictions and then just parse from that.
func Write_new_prediction_outcomes(db *sqlite3.Conn, predictions []map[string]interface{}) error {
	
	for i := 0; i < len(predictions); i++{
		prediction_data := predictions[i]
		prediction_id := prediction_data["prediction_id"]
		outcome_id := prediction_data["outcome_id"]
		title := prediction_data["title"]
		lose_win := prediction_data["lose_win"]
		sql_query_string := fmt.Sprintf(`INSERT INTO outcomes ('prediction_id', 'outcome_id', 'title', 'lose_win') VALUES('%s', '%s','%s', %d)`,prediction_id, outcome_id, title, lose_win)	
		err := db.Exec(sql_query_string)

		if err!=nil{
			err = errors.New("FILE: sqlite_helper FUNC: Write_new_prediction_outcomes CALL: db.Exec " + err.Error())
			return err
		}
	}
	
	return nil
}

//This returns the prediction id, prediction created at, and an error.
func Get_predictions(db *sqlite3.Conn, sub string, status string) (string,string, error) {
	
	sql_query_string := fmt.Sprintf(`SELECT * FROM prediction WHERE broadcaster_id == '%s' AND status == '%s'`, sub, status)
	sql_query, _, err := db.Prepare(sql_query_string)
	if err != nil{
		err = errors.New("FILE: sqlite_helper FUNC: Get_predictions CALL: db.Prepare " + err.Error())
		return "", "" ,err
	}
	prediction_id := ""
	created_at := ""
	for sql_query.Step() {
		prediction_id = sql_query.ColumnText(1)
		created_at = sql_query.ColumnText(3)
	}
	if prediction_id == ""{
		//err = errors.New("FILE: sqlite_helper FUNC: Get_predictions BUG: prediction_id was blank")
		return "null", "null", nil
	}
	// if created_at == ""{
	// 	//err = errors.New("FILE: sqlite_helper FUNC: Get_predictions BUG: created_at was blank")
	// 	return "null", "null", nil
	// }
	return prediction_id, created_at, nil	
}

func Get_prediction_outcome_id(db *sqlite3.Conn, prediction_id string, lose_win int)(string, error){
	
	sql_query_string := fmt.Sprintf(`SELECT * FROM outcomes WHERE prediction_id == '%s' AND lose_win == '%d'`, prediction_id, lose_win)
	sql_query, _, err := db.Prepare(sql_query_string)
	if err != nil{
		err = errors.New("FILE: sqlite_helper FUNC: Get_predictions_outcome_id CALL: db.Prepare " + err.Error())
		return "", err
	}
	outcome_id := ""
	for sql_query.Step() {
		outcome_id = sql_query.ColumnText(1)
	} 
	if outcome_id ==""{
		err = errors.New("FILE: sqlite_helper FUNC: Get_predictions BUG: outcome_id was blank")
		return "", err
	}

	return outcome_id, nil		
}

//This first deletes the prediction outcomes associated with the prediction, then deletes the prediction id.
func Delete_prediction_id(db *sqlite3.Conn, sub string)error{
	prediction_id, _, err := Get_predictions(db, sub, "ACTIVE")
	if err!=nil{
		return err
	}
	err = Delete_outcomes(db, prediction_id)
	if err!=nil{
		return err
	}
	sql_query_string := fmt.Sprintf(`DELETE FROM prediction WHERE broadcaster_id == '%s'`, prediction_id)
	err = db.Exec(sql_query_string)
	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Delete_prediction_id CALL: db.Exec " + err.Error())
		return err
	}
	return nil
}

func Delete_outcomes(db *sqlite3.Conn, prediction_id string) error{
	sql_query_string := fmt.Sprintf(`DELETE FROM outcomes WHERE prediction_id == '%s'`, prediction_id)
	err := db.Exec(sql_query_string)
	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Delete_outcomes CALL: db.Exec " + err.Error())
		return err
	}
	return err
}
func Delete_all_predictions(db *sqlite3.Conn, sub string) error{
	
	sql_query_string := fmt.Sprintf(`DELETE FROM prediction WHERE broadcaster_id == '%s'`, sub)
	err := db.Exec(sql_query_string)
	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Delete_all_predictions CALL: db.Exec " + err.Error())
		return err
	}
	return err
}

func Write_sub_event(db *sqlite3.Conn, event_id string) error{

	sql_query_string := fmt.Sprintf(`INSERT INTO Sub_Events (Sub_Event_ID) VALUES ('%s')`, event_id)
	err := db.Exec(sql_query_string)
	if err != nil{
		err = errors.New("FILE: sqlite_helper FUNC: Write_sub_event CALL: db.Exec " + err.Error())
		return err
	}
	return nil
}

func Get_sub_event(db *sqlite3.Conn, event_id string)(bool, error){
	sql_query_string := fmt.Sprintf(`SELECT * FROM Sub_Events WHERE Sub_Event_ID == '%s'`, event_id)

	sql_query, _, err := db.Prepare(sql_query_string)

	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Get_sub_event CALL: db.Prepare " + err.Error())
		return false, err
	}

	if sql_query.Step() {
		return true, nil
	} 

	return false, nil
}

func Get_all_access_tokens(db *sqlite3.Conn)([]Twitch_user_refresh, error){

	sql_query_string := `SELECT * FROM twitch_user_info`

	var Refresh_list []Twitch_user_refresh
	var Twitch_refresh Twitch_user_refresh

	sql_query, _ ,err := db.Prepare(sql_query_string)

	if err !=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Get_all_access_tokens CALL: db.Prepare " + err.Error())
		return Refresh_list, err
	}

	for sql_query.Step(){
		Twitch_refresh.Access_token = sql_query.ColumnText(2)
		Twitch_refresh.Refresh_token = sql_query.ColumnText(3)
		Twitch_refresh.User_id = sql_query.ColumnText(0)
		Refresh_list = append(Refresh_list, Twitch_refresh)
	}

	return Refresh_list, nil
}

func Update_state(db *sqlite3.Conn, state string, player_tag string)(error){

	sql_query_string := fmt.Sprintf(`UPDATE state SET player_id = '%s' WHERE state_value == '%s'`,player_tag, state)
	err := db.Exec(sql_query_string)

	if err!=nil{
		err = errors.New("FILE: sqlite_helper FUNC: Update_state CALL: db.Exec " + err.Error())
		return err}

	return nil
}