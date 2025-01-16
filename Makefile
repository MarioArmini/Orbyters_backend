swag:
	swag init --parseDependency --parseInternal -g main.go

run:
	go run main.go

.PHONY: swag run