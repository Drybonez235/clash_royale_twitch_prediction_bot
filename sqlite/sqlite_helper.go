package sqlite

import (
	"errors"
	"fmt"

	//"github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const file = "file:twitch_authorization"

type Twitch_user struct{
	//sql_table_values := "'sub', 'display_name',
	// 'access_token', 'refresh_token', 'scope', 
	//'token_type', 'app_request', 'app_received', 'token_exp', 'token_iat', 'token_iss'"
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
}

func open_db()(*sqlite3.Conn, error){
	db, err := sqlite3.Open(file) 

	if err != nil {
		err = errors.New("there was a problem opening the database file")
		return db, err
	}	

	return db, nil
}

func Create_twitch_database() error {

	db, err := open_db() 

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Exec(`CREATE TABLE state (state_value text)`)
	if err != nil{
		err = errors.New("there was a problem creating the  state table")
		return err
	}

	err = db.Exec(`CREATE TABLE twitch_user_info 
	(sub text, display_name text, access_token text, refresh_token text, scope text, token_type text, app_request text,
	app_received text, token_exp float, token_iat float, token_iss text)`)

	if err != nil{
		err = errors.New("there was a problem creating the twitch_user_info table")
		return err
	}

	return err
}

func Write_state_nonce(state_nonce string, table string) error {

	//fmt.Println("Wrote user to state or nonce")
	db, err := open_db()

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return err
	}

	header := ""

	if table == "state"{
		header = "state_value"
	} else if table == "nonce" {
		header = "nonce_value"
	} else {
		err = errors.New("invalid table given")
	}

	if err != nil{
		return err
	}

	sql_command := fmt.Sprintf(`INSERT INTO '%s' ('%s') VALUES ('%s')`, table, header, state_nonce)
	fmt.Println(sql_command)

	err = db.Exec(sql_command) //`INSERT INTO state (state_value) VALUES ('Testing')`

	if err != nil{ 
		err = errors.New("there was a problem inserting into state")
		return err
	}

	err = db.Close()

	return err
}

func Check_state_nonce(state_nonce string, table string) (bool, error){
	fmt.Println("Checked for state or nonce")

	db, err := open_db()

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return  false, err
	}

	header := ""

	if table == "state"{
		header = "state_value"
	} else if table == "nonce" {
		header = "nonce_value"
	} else {
		err = errors.New("invalid table given")
	}

	if err != nil{
		return false, err
	}

	sql_query_string := fmt.Sprintf(`SELECT * FROM '%s' WHERE %s == '%s'`, table, header, state_nonce)

	sql_query, _, err := db.Prepare(sql_query_string)

	if err != nil{
		err = errors.New("there was a problem prepairing the query")
		return false ,err
	}

	if sql_query.Step() {
		sql_query.Close()
		err = delete_state_nonce(state_nonce, table, db)
		return true, err
	} else {
		sql_query.Close()
	}

	return false, err
}

func delete_state_nonce(state_nonce string, table string, db *sqlite3.Conn) error {
	fmt.Println("Deleted state or Nonce")

	header := ""

	if table == "state"{
		header = "state_value"
	} else if table == "nonce" {
		header = "nonce_value"
	} else {
		err := errors.New("invalid table given")
		return err
	}

	sql_query_string := fmt.Sprintf(`DELETE FROM '%s' WHERE '%s' == '%s'`, table, header, state_nonce )

	fmt.Println(sql_query_string)

	err := db.Exec(sql_query_string)

	//fmt.Println(err.Error())

	if err!=nil{
		err = errors.New("there was a problem deleting the record")
		return err
	}

	err = db.Close()

	return err
}

func Write_twitch_info(sub string, display_name string, access_token string, refresh_token string, scope string, token_type string,
	 app_request string, app_received string, token_exp int, token_iat int, token_iss string) error {

		err := Remove_twitch_user(sub)

		if err!=nil{
			return err
		}

		fmt.Println("Wrote user to Dtabase")

		db, err := open_db()

		if err!=nil{
			err = errors.New("there was a problem opening the database")
			return  err
		}

		//(sub text, access_token text, refresh_token text, scope text, token_type text, app_request text
			//app_received text, token_exp float, token_iat float, token_iss text)
		sql_table_values := "'sub', 'display_name', 'access_token', 'refresh_token', 'scope', 'token_type', 'app_request', 'app_received', 'token_exp', 'token_iat', 'token_iss'"

		sql_command := fmt.Sprintf("INSERT INTO twitch_user_info (%s) VALUES ('%s', '%s' , '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s')", sql_table_values, sub, display_name, access_token, refresh_token, scope, token_type, app_request, app_received, token_exp, token_iat, token_iss)
		
		fmt.Println(sql_command)

		err = db.Exec(sql_command)	

		if err!=nil{
			fmt.Println(err)
			err = errors.New("there was a problem inserting the twitch user")
		}

		return err
}

func Get_twitch_user(id_type string, id string) (Twitch_user, error){
	var twitch_user Twitch_user

	fmt.Println("Retrieved twitch infor from db")
	db, err := open_db()

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return  twitch_user, err
	}

	field := ""

	if id_type == "sub"{
		field = "sub"
	} else if id_type == "display_name" {
		field = "display_name"
	} else {
		err = errors.New("invalid id type")
		return twitch_user, err
	}
 
	sql_query_string := fmt.Sprintf(`SELECT * FROM twitch_user_info WHERE %s == '%s'`, field, id)

	sql_query, _, err := db.Prepare(sql_query_string)

	if err != nil{
		err = errors.New("there was a problem prepairing the query")
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
	} 

	defer db.Close()

	return twitch_user, err	
}

func Update_tokens(access_token string, refresh_token string, sub string)(error){
	db, err := open_db()

	if err!=nil{
		return err
	}
	defer db.Close()

	sql_query_string := fmt.Sprintf(`UPDATE twitch_user_info SET access_token = '%s', refresh_token = '%s' WHERE sub=='%s'`, access_token, refresh_token, sub)

	err = db.Exec(sql_query_string)

	if err!=nil{
		fmt.Println("We got here right?")
		fmt.Println(err)
		return err
	}
	fmt.Println(sql_query_string)

	fmt.Println("Token was updated")
	return nil
}

func Remove_twitch_user(sub string) error {
	db, err := open_db()

	if err!=nil{
		return err
	}

	sql_string := fmt.Sprintf("DELETE FROM twitch_user_info WHERE sub=='%s'", sub)

	err = db.Exec(sql_string)

	if err!=nil{
		return err
	}

	return nil
}