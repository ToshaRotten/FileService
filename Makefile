.PHONY: serv cli
build_server:
	cd server/ && go build main.go
build_client:docker build -t go-docker-image .
	cd client/ && go build main.go

server: build_server
	cd server/ && ./main
client: build_client
	cd client/ && ./main
