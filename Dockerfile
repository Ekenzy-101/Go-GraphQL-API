# golang:1.16-buster 
FROM golang@sha256:9cb4d9fe93e7efd88048e2a30037b38547e0c61ff171e381b8a499b189a71817 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux go build -o main .

FROM gcr.io/distroless/base

COPY --from=builder /app/main /

EXPOSE 5000

ENTRYPOINT ["/main"]