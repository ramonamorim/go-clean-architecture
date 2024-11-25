FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

FROM scratch

COPY --from=builder /app/main main
COPY --from=builder /app/.env .env

EXPOSE 8000

CMD ["/main"]