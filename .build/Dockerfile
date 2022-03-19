FROM golang:1.17

RUN mkdir "build"

COPY ../ /build
WORKDIR /build

RUN GOOS=linux CGO_ENABLED=0 go build -o ./app ./cmd/server/main.go

FROM alpine:latest
COPY --from=0 /.build/app ./
COPY ../config/config.yaml ./config/config.yaml

CMD ["./app"]