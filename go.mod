module github.com/Drybonez235/clash_royale_twitch_prediction_bot

go 1.22.3

//replace github.com/Drybonez235/clash_royale_twitch_prediction_bot/twitch_api => /Drybonez235/clash_royale_twitch_prediction_bot/twitch_api

//replace github.com/Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api => /Drybonez235/clash_royale_twitch_prediction_bot/clash_royale_api

require github.com/ncruces/go-sqlite3 v0.22.0

require (
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.8.2 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
