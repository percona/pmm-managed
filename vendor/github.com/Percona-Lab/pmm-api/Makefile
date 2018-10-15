all: gen

gen:
	go install -v ./vendor/github.com/golang/protobuf/protoc-gen-go
	go run generate.go
	go install -v ./...
