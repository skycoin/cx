
build:
	go build -o packageloader

test_loader:
	go test ./loader -v -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

test_encoder:
	go run main.go -load -name Test -path ./loader/testValidProgram/
	go test ./encoder -v -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html
	sudo rm -r ./encoder/test/