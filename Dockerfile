FROM golang:1.26
ADD . /app
WORKDIR /app
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o nuki-cron cmd/*.go
WORKDIR /
ENTRYPOINT ["/app/nuki-cron"]