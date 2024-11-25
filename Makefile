up:
	docker compose up -d

down:
	docker compose down

run-service:
	cd cmd/ordersystem && go run main.go wire_gen.go

