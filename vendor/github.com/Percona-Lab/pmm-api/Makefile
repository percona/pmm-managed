all: gen

init:
	# https://github.com/uber/prototool#installation
	curl -L https://github.com/uber/prototool/releases/download/v1.3.0/prototool-$(shell uname -s)-$(shell uname -m) -o ./prototool
	chmod +x ./prototool

clean:
	find . -name '*.pb.go' -not -path './vendor/*' -delete
	find . -name '*.pb.gw.go' -not -path './vendor/*' -delete
	find . -name '*.swagger.json' -not -path './vendor/*' -delete

	rm -fr http

gen: clean
	go install -v ./vendor/github.com/golang/protobuf/protoc-gen-go \
					./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
					./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
					./vendor/github.com/go-swagger/go-swagger/cmd/swagger

	./prototool all

	# no public API for pmm-agent
	rm -f agent/*.swagger.json

	swagger mixin inventory/inventory.json inventory/*.swagger.json --output=inventory.swagger.json
	swagger validate inventory.swagger.json
	rm -f inventory/*.swagger.json

	mkdir http
	swagger generate client --spec=inventory.swagger.json --target=http
	go install -v ./...

serve:
	go run swagger/serve.go -dir=swagger
