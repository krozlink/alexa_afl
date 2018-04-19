build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/win_next functions/win_next/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/win_premiership functions/win_premiership/main.go
	env go build -ldflags="-s -w" -o bin/test main.go
	