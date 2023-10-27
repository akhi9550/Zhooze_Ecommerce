run:
	go run main.go


swag: ## Generate swagger docs
	swag init -g main.go -o ./cmd/docs