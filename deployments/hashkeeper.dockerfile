FROM golang:1.21.7 as builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN GO111MODULE=auto CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -o ./bin/hashkeeper ./hashkeeper/cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/hashkeeper .

ENTRYPOINT ["./hashkeeper"]
