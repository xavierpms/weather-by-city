FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -C "cmd/server" -ldflags="-w -s" -o weather-by-city .

FROM scratch
COPY --from=builder /app/cmd/server/weather-by-city .
CMD ["./weather-by-city"]
