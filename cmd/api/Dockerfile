FROM golang:1.23.1-alpine AS build

ENV GOPATH="" CGO_ENABLED=0

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/api cmd/api
COPY cmd/api/config.json /config.json
COPY internal internal

RUN go mod download
RUN go build -o /bin/app cmd/api/main.go

FROM alpine:3.20.1

COPY --from=build /bin/app /app/
COPY --from=build /config.json /app/config.json

WORKDIR /app/
EXPOSE 8080

CMD ["./app", "-c=config.json"]