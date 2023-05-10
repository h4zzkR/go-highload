go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$ export PATH="$PATH:$(go env GOPATH)/bin"
$ go get golang.org/x/net/context

$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/messenger.proto

$ sudo snap install redis
$ go get github.com/redis/go-redis/v9

# Run
run redis server: `redis-server`
server: `go build && ./grpc-messenger -s  `
client with h4zzkR username: `go build && ./grpc-messenger -u h4zzkR`

# Demo

People talk in spaces (or telegram groups), there is no direct messages.
Messages are server-side cached, when new client connects he loads this message history and can view older messages.

![two-people-session](https://lh3.googleusercontent.com/pw/AMWts8DOruQBrr2YttzHXzIlaqLvfxugsAOaq9XnmpayfGB1tlnPZeM5vixDWExeIRqlS09faq-C-eHzFm_iP1ZpDCglWykz22EN1GlcVELrKP-vCHK76gqn0gCNfpBZ52F7kbkw7ssnvF2A_4VCosSTzzGO=w1280-h313-s-no?authuser=0)

# Todo
- Cancellation (shutdown server, logout, client reactions to this events)
