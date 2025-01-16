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
		err = errors.New("there was a problem creating the table")
		return err
	}

	return err
}

func Write_state(state string) error {
	db, err := sqlite3.Open(file)

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return err
	}

	sql_command := fmt.Sprintf(`INSERT INTO state (state_value) VALUES ('%s')`, state)
	fmt.Println(sql_command)

	err = db.Exec(sql_command) //`INSERT INTO state (state_value) VALUES ('Testing')`

	if err != nil{ 
		err = errors.New("there was a problem inserting into state")
		return err
	}

	err = db.Close()

	return err
}

func Check_state(state string) (bool, error){
	db, err := sqlite3.Open(file)

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return  false, err
	}

	sql_query_string := fmt.Sprintf(`SELECT * FROM state WHERE state_value == '%s'`, state)

	sql_query, _, err := db.Prepare(sql_query_string)

	if err != nil{
		err = errors.New("there was a problem prepairing the query")
		return false ,err
	}

	defer sql_query.Close()

	if sql_query.Step() {
		err = delete_state(state)
		return true, err
	} else {
		delete_state(state)
	}

	return false, err
}

func delete_state(state string) error {
	db, err := sqlite3.Open(file)

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return err
	}

	sql_query_string := fmt.Sprintf(`DELETE FROM state WHERE state_value == '%s'`, state)

	err = db.Exec(sql_query_string)

	if err!=nil{
		err = errors.New("there was a problem deleting the record")
		return err
	}

	err = db.Close()

	return err
}
