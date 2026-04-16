FROM golang:1.26-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /app/bin/symlinker .
