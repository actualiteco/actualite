FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o actualite cmd/actualite/main.go

FROM alpine:3.23 AS final

WORKDIR /app

COPY --from=builder /build/actualite .

COPY --from=builder /build/.env .

CMD [ "/app/actualite" ]