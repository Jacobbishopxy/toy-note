
dev:
	cd cmd/app && go run .

build:
	cd cmd/app && go build .

prod:
	cd cmd/app && ./app -m prod

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

swag-fmt:
	swag fmt -d ./api -g ../cmd/app/main.go

swag-init: swag-fmt
	swag init -d ./api -g ../cmd/app/main.go
