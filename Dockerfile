FROM golang:latest

RUN go install github.com/ToshaRotten/fileService@latest

WORKDIR /app

COPY ./ /app

RUN make server
RUN make client