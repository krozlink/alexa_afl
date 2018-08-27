build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/win_next functions/win_next/main.go
	