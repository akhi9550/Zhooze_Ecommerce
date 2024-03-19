run:
	go run main.go


swag: 
	swag init -g main.go -o ./cmd/docs

# go get -u github.com/swaggo/swag/cmd/swag
# swag init
# go install github.com/swaggo/swag/cmd/swag@latest
