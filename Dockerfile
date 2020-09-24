FROM golang:alpine as builder

RUN pwd
WORKDIR /go/src/
COPY . /go/src/

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/tictactoe main.go

# Prepare final, minimal image
FROM alpine:latest
WORKDIR /app
RUN mkdir -p /app/static/

COPY --from=builder /go/src/bin/ /app
COPY --from=builder /go/src/static /app/static

# Install
ENV HOME /app

CMD ["./tictactoe"]