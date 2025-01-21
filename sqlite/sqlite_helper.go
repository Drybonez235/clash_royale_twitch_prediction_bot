package sqlite

import (
	//"database/sql"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"errors"
	"fmt"
)

const file = "file:twitch_authorization"

func Create_twitch_database() error {
	db, err := sqlite3.Open(file) 

	if err != nil {
		err = errors.New("there was a problem opening the database file")
		return err
	}

	err = db.Exec(`CREATE TABLE state (state_value text)`)
	if err != nil{
		err = errors.New("there was a problem creating the  state table")
		return err
	}

	err = db.Exec(`CREATE TABLE nonce (nonce_value text)`)

	if err != nil{
		err = errors.New("there was a problem creating the nance table")
		return err
	}

	return err
}

func Write_state_nonce(state_nonce string, table string) error {
	db, err := sqlite3.Open(file)

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
	db, err := sqlite3.Open(file)

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
		err = delete_state_nonce(state_nonce, table)
		return true, err
	} else {
		sql_query.Close()
	}

	return false, err
}

func delete_state_nonce(state_nonce string, table string) error {
	db, err := sqlite3.Open(file)

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

	sql_query_string := fmt.Sprintf(`DELETE FROM '%s' WHERE '%s' == '%s'`, table, header, state_nonce )

	fmt.Println(sql_query_string)

	err = db.Exec(sql_query_string)

	//fmt.Println(err.Error())

	if err!=nil{
		err = errors.New("there was a problem deleting the record")
		return err
	}

	err = db.Close()

	return err
}
