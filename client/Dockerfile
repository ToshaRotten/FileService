FROM golang:latest

WORKDIR /client
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

EXPOSE 9090

RUN go build client.go

CMD [ "./client" ]
#docker run --rm -it --name prg fileservice-client