swag:
	swag init -g cmd/server/main.go

run:
	cd cmd/server && go run main.go

.PHONY: swag run