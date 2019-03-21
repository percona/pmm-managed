all:
	go install -v ./...
	go test -c -v ./inventory
	go test -c -v ./server

race:
	go install -v -race ./...
	go test -c -v -race ./inventory
	go test -c -v -race ./server
