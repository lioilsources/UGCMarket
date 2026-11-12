# market-api: multi-stage, the result is a single static binary (modeled on ugc-api).
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /market-api ./cmd/market-api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /market-api /usr/local/bin/market-api
VOLUME /data
EXPOSE 8096
ENTRYPOINT ["market-api"]
