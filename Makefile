run:
	export CLIENT_URL=http://localhost:3000/
	go run ./cmd/RPC/main.go

test:
	go test ./tests