# agonake-server
snake server using agones

## instructions for development
1. Run agones sdk server: `./agonessdk-server-0.1/sdk-server.<platform>.amd64 --local`.
    special note: Windows binary needs `.exe` at the end of the file name.
2. Run go program as usually: `go run cmd/agonake/main.go`.
3. For testing commands through terminal: `nc -u localhost 7654`.

## instructions for kubernetes
1. Build go binary with: `CGO_ENABLED=0 GOOS=linux go build -o ./bin/agonake-server cmd/agonake-server/main.go`.
2. Build go image with: `docker build . --tag=jmaeso/agonake-server:0.1`.
