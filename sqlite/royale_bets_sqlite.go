package sqlite

import(
	"fmt"
	"errors"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Royale_bets_viewer struct{
	Session_id int
	Streamer_player_tag string
	Screen_name string
	Total_points int
	Last_refresh_time int
}

type Royale_bets_streamer struct{
	Stream_start_time int `json:"stream_start_time"`
	Streamer_player_tag string `json:"streamer_player_tag"`
	Streamer_last_refresh_time int	`json:"streamer_last_refresh_time"`
	Wins int `json:"wins"`
	Losses int 	`json:"losses"`
}

type Battle_result struct{
	Player_tag string `json:"streamer_player_tag"`
	Battle_time int `json:"battle_time"`
	Red_crowns_taken int `json:"crowns_taken_int"`
	Blue_crowns_lost int `json:"crowns_lost_int"`
}

type Leader_board_entry struct{
	Rank int `json:"rank"`
	Screen_name string `json:"screen_name"`
	Total_points int `json:"total_points"`
}

func Create_royale_bets_db(db *sqlite3.Conn) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS RoyaleBetsViewer (
			Session_id INTEGER PRIMARY KEY,
			Streamer_player_tag TEXT,
			Screen_name TEXT,
			Total_points INTEGER,
			Last_refresh_time INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS RoyaleBetsStreamer (
			Stream_start_time INTEGER PRIMARY KEY,
			Streamer_player_tag TEXT,
			Streamer_last_refresh_time INTEGER,
			Wins INTEGER,
			Losses INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS BattleResult (
			Player_tag TEXT,
			Battle_time INTEGER,
			Red_crowns_taken INTEGER,
			Blue_crowns_lost INTEGER
		);`,
	}//PRIMARY KEY (Player_tag, Battle_time)

	for _, query := range queries {
		if err := db.Exec(query); err != nil {
			return errors.New("Error creating tables: " + err.Error())
		}
	}

	return nil
}

func Insert_royale_bets_viewer(db *sqlite3.Conn, viewer Royale_bets_viewer) error {
	stmt, _, err := db.Prepare("INSERT INTO RoyaleBetsViewer (Session_id, Streamer_player_tag, Screen_name, Total_points, Last_refresh_time) VALUES (?, ?, ?, ?,?)")
	if err != nil {
		return errors.New("Error preparing statement for RoyaleBetsViewer: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindInt(1, viewer.Session_id); err != nil {
		return err
	}
	if err := stmt.BindText(2, viewer.Streamer_player_tag); err != nil {
		return err
	}
	if err := stmt.BindText(3, viewer.Screen_name); err != nil {
		return err
	}
	if err := stmt.BindInt(4, viewer.Total_points); err != nil {
		return err
	}
	if err := stmt.BindInt(5, viewer.Last_refresh_time); err != nil{
		return err
	}
	return stmt.Exec()
}

func Insert_royale_bets_streamer(db *sqlite3.Conn, streamer Royale_bets_streamer) error {
	stmt, _, err := db.Prepare("INSERT INTO RoyaleBetsStreamer (Stream_start_time, Streamer_player_tag, Streamer_last_refresh_time, Wins, Losses) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error preparing statement for RoyaleBetsStreamer: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindInt(1, streamer.Stream_start_time); err != nil {
		return err
	}
	if err := stmt.BindText(2, streamer.Streamer_player_tag); err != nil {
		return err
	}
	if err := stmt.BindInt(3, streamer.Streamer_last_refresh_time); err != nil {
		return err
	}
	if err := stmt.BindInt(4, streamer.Wins); err != nil {
		return err
	}
	if err := stmt.BindInt(5, streamer.Losses); err != nil {
		return err
	}
	return stmt.Exec()
}

func Insert_battle_result(db *sqlite3.Conn, result Battle_result) error {
	fmt.Println("Battle result added to db: ", result)
	stmt, _, err := db.Prepare("INSERT INTO BattleResult (Player_tag, Battle_time, Red_crowns_taken, Blue_crowns_lost) VALUES (?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error preparing statement for BattleResult: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindText(1, result.Player_tag); err != nil {
		return err
	}
	if err := stmt.BindInt(2, result.Battle_time); err != nil {
		return err
	}
	if err := stmt.BindInt(3, result.Red_crowns_taken); err != nil {
		return err
	}
	if err := stmt.BindInt(4, result.Blue_crowns_lost); err != nil {
		return err
	}

	return stmt.Exec()
}

func Get_royale_bets_viewer(db *sqlite3.Conn, session_id int, screen_name string) (*Royale_bets_viewer, error) {
	stmt, _, err := db.Prepare("SELECT Session_id, Streamer_player_tag, Screen_name, Total_points, Last_refresh_time FROM RoyaleBetsViewer WHERE Session_id = ? AND Screen_name = ?")
	if err != nil {
		return nil, errors.New("Error preparing statement for get_royale_bets_viewer: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindInt(1, session_id); err != nil {
		return nil, err
	}
	if err := stmt.BindText(2, screen_name); err != nil {
		return nil, err
	}

	var viewer Royale_bets_viewer
	if stmt.Step() {
		viewer.Session_id = stmt.ColumnInt(0)
		viewer.Streamer_player_tag = stmt.ColumnText(1)
		viewer.Screen_name = stmt.ColumnText(2)
		viewer.Total_points = stmt.ColumnInt(3)
		viewer.Last_refresh_time = stmt.ColumnInt(4)
		return &viewer, nil
	}

	return nil, errors.New("No viewer found")
}

func Get_royale_bets_streamer(db *sqlite3.Conn, streamer_player_tag string, viewer_session_id int) (*Royale_bets_streamer, error) {

	stmt, _, err := db.Prepare("SELECT Stream_start_time, Streamer_player_tag, Streamer_last_refresh_time, Wins, Losses FROM RoyaleBetsStreamer WHERE Streamer_player_tag = ? AND Stream_start_time > (? - (60000 * 60 * 12))")
	if err != nil {
		return nil, errors.New("Error preparing statement for RoyaleBetsStreamer: " + err.Error())
	}
	defer stmt.Close()

	// Bind the Streamer_player_tag parameter
	if err := stmt.BindText(1, streamer_player_tag); err != nil {
		return nil, err
	}

	if err := stmt.BindInt(2, viewer_session_id); err!=nil{
		return nil, err
	}

	// Execute the query and fetch data
	if stmt.Step() {
		streamer := Royale_bets_streamer{
			Stream_start_time:         stmt.ColumnInt(0),
			Streamer_player_tag:       stmt.ColumnText(1),
			Streamer_last_refresh_time: stmt.ColumnInt(2),
			Wins:                      stmt.ColumnInt(3),
			Losses:                    stmt.ColumnInt(4),
		}
		return &streamer, nil
	}
	
	return nil, nil
}

func Update_royale_bets_viewer(db *sqlite3.Conn, session_id int, screen_name string, total_points int, last_refresh_time int) error {

	stmt, _, err := db.Prepare("UPDATE RoyaleBetsViewer SET Total_points = ?, Last_refresh_time = ? WHERE Session_id = ? AND Screen_name = ?")
	if err != nil {
		return errors.New("Error preparing statement for update_royale_bets_viewer: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindInt(1, total_points); err != nil {
		return err
	}
	if err := stmt.BindInt(2, last_refresh_time); err != nil {
		return err
	}
	if err := stmt.BindInt(3, session_id); err != nil {
		return err
	}
	if err := stmt.BindText(4, screen_name); err != nil {
		return err
	}
	fmt.Println("Updated viewer: Total Points: ", total_points, " Last Refresh Time:", last_refresh_time )
	return stmt.Exec()
}

func Update_royale_bets_streamer_wins_losses(db *sqlite3.Conn, player_tag string, stream_start_time int , last_refresh_time int, win_lose string)(error){
	fmt.Println("Update Streamer: Last Refresh Time ", last_refresh_time)
	
	var column string
	switch win_lose {
	case "win": 
		column = " Wins = Wins + 1"
	case "lose":
		column = "Losses = Losses + 1"
	} 

	query := "UPDATE RoyaleBetsStreamer SET " + column + ", Streamer_last_refresh_time = ? WHERE Streamer_player_tag = ? AND Stream_start_time = ?"
	
	stmt, _, err := db.Prepare(query)
	if err!=nil{
		return errors.New("FILE: Royale_bets_sqlite FUNC: Update_royale_bets_streamer_wins_losses CALL: db.prepare " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindText(2, player_tag); err!= nil{
		return err
	}

	if err := stmt.BindInt(3, stream_start_time); err!= nil{
		return err
	}

	if err:= stmt.BindInt(1, last_refresh_time); err!= nil{
		return err
	}
//I WASNT ACTULLY CARRING OUT THE SQL
	return stmt.Exec()
}

func Get_battle_result(db *sqlite3.Conn, player_tag string, last_refresh_time int) ([]Battle_result, error) {
	stmt, _, err := db.Prepare("SELECT Player_tag, Battle_time, Red_crowns_taken, Blue_crowns_lost FROM BattleResult WHERE Player_tag = ? AND Battle_time >= ?")
	if err != nil {
		return nil, errors.New("Error preparing statement for get_battle_result: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindText(1, player_tag); err != nil {
		return nil, err
	}
	if err := stmt.BindInt(2, last_refresh_time); err != nil {
		return nil, err
	}

	results  := []Battle_result{}
	for stmt.Step() {
		result := Battle_result{
			Player_tag:      stmt.ColumnText(0),
			Battle_time:     stmt.ColumnInt(1),
			Red_crowns_taken: stmt.ColumnInt(2),
			Blue_crowns_lost: stmt.ColumnInt(3),
		}
		results = append(results, result)
	}
	return results, nil
}
func Get_all_battle_results(db *sqlite3.Conn) ([]Battle_result, error) {
	stmt, _, err := db.Prepare("SELECT * FROM BattleResult")
	if err != nil {
		return nil, errors.New("Error preparing statement for get_battle_result: " + err.Error())
	}
	defer stmt.Close()

	results  := []Battle_result{}
	for stmt.Step() {
		result := Battle_result{
			Player_tag:      stmt.ColumnText(0),
			Battle_time:     stmt.ColumnInt(1),
			Red_crowns_taken: stmt.ColumnInt(2),
			Blue_crowns_lost: stmt.ColumnInt(3),
		}
		results = append(results, result)
	}
	return results, nil
}

func Get_top_ten_and_viewer_position(db *sqlite3.Conn, player_tag string, stream_start_time int,viewer_session_id int) (*[]Leader_board_entry, error) {
	stmt, _, err := db.Prepare("SELECT Screen_name, Total_points, Session_id FROM RoyaleBetsViewer WHERE Streamer_player_tag = ? AND Session_id >= ? ORDER BY Total_points")
	if err != nil {
		return nil, errors.New("Error preparing statement for get_battle_result: " + err.Error())
	}
	defer stmt.Close()

	if err := stmt.BindText(1, player_tag); err != nil {
		return nil, err
	}
	if err := stmt.BindInt(2, stream_start_time); err != nil {
		return nil, err
	}

	rank := 0
	found := false
	var viewer_rank Leader_board_entry

	var entries []Leader_board_entry

	for stmt.Step() {
		var entry Leader_board_entry
		if rank < 10 {
			entry.Rank = rank + 1
			entry.Screen_name = stmt.ColumnText(0)
			entry.Total_points = stmt.ColumnInt(1)
			entries = append(entries, entry)
		} 
		if stmt.ColumnInt(2) == viewer_session_id {
			viewer_rank.Rank = rank + 1
			viewer_rank.Screen_name = stmt.ColumnText(0)
			viewer_rank.Total_points = stmt.ColumnInt(1)
			found = true
		}
		rank += 1
	}

	if found{
		entries = append(entries, viewer_rank)
	}
	return &entries, nil
}
