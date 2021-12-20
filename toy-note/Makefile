
dev:
	cd cmd/app && go run .

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

swag-fmt:
	swag fmt -d ./api -g ../cmd/app/main.go

swag-init:
	swag init -d ./api -g ../cmd/app/main.go

