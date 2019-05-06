# simple-chat-backend

## How to run
1. Have a MySQL server and database running on your machine (or in another one), and write the credentials to the config.json file
2. Have a RabbitMQ server running on your machine (or in another one), and write the credentials to the config.json file
3. `go run *.go` or if you are running it on a Windows machine, `go run auth.go bot.go commands.go config.go database.go rooms.go router.go server.go users.go`

## How to test
1. `go test -v`