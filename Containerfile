FROM docker.io/golang:1.24.4 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

FROM docker.io/alpine:3.22.0
COPY --from=builder /app/app /bin/app
ENTRYPOINT ["/bin/app"]
