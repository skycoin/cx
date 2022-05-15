
build:
	go build -o packageloader

test_loader:
	go test ./loader -v -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

test_server:
	sudo systemctl restart redis.service
	go run main.go -l Test2 ./loader/test2
	redis-cli