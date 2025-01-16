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

	err = db.Exec(`INSERT INTO test VALUES (0), (1), (2)`)

	if err != nil{
		err = errors.New("there was a problem insterting into the database")
		return err
	}

	prepaired_query, _, err := db.Prepare(`SELECT * FROM test`)

	if err != nil{
		err = errors.New("there was a problem prepairing the query")
		return err
	}

	defer prepaired_query.Close()

	for prepaired_query.Step() {
		fmt.Println(prepaired_query.ColumnInt(0))
	}

	err = db.Close()

	return err
}

func write_state(state string) error {
	db, err := sqlite3.Open(file)

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return err
	}

	sql_command := fmt.Sprintf("INSERT INTO secrets VALUES(%s)", state)

	err = db.Exec(sql_command)

	if err != nil{
		err = errors.New("there was a porblem insterting into secrets")
		return err
	}

	err = db.Close()

	return err
}

func check_state(state string) (bool, error){
	db, err := sqlite3.Open(file)

	if err!=nil{
		err = errors.New("there was a problem opening the database")
		return  false, err
	}

	sql_query_string := fmt.Sprintf("SELECT 1 FROM state WHERE state_value == %s", state)

	sql_query, _, err := db.Prepare(sql_query_string)

	if err != nil{
		err = errors.New("there was a problem prepairing the query")
		return false ,err
	}

	defer sql_query.Close()

	if sql_query.Step() {
		delete_state(state)
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

	sql_query_string := fmt.Sprintf("DELETE FROM state WHERE state_value == %s", state)

	err = db.Exec(sql_query_string)

	if err!=nil{
		err = errors.New("there was a problem deleting the record")
		return err
	}

	err = db.Close()

	return err
}
