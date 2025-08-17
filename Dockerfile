# golang:1.24 
FROM golang@sha256:e155b5162f701b7ab2e6e7ea51cec1e5f6deffb9ab1b295cf7a697e81069b050 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux go build -o main .

FROM gcr.io/distroless/base

COPY --from=builder /app/main /

EXPOSE 5000

ENTRYPOINT ["/main"]