module github.com/thekhanj/csdmpro

go 1.24.0

require (
	github.com/PuerkitoBio/goquery v1.10.2
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/google/wire v0.6.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/thekhanj/tgool v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.38.0
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/thekhanj/drouter v0.0.1 // indirect
)

replace github.com/thekhanj/tgool => ../tgool
