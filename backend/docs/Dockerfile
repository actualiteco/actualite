FROM golang:1.27-alpine

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /code

ENTRYPOINT ["swag"]
