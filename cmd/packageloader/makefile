
build:
	go build -o packageloader

test_loader:
	go test ./loader -v -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

test_encoder:
	make build
	sudo rm -r ./encoder/test_redis/
	sudo rm -r ./encoder/test_bolt/
	./packageloader -redis -load -name Test -path ./loader/test_folder/test_valid_program
	./packageloader -load -name Test -path ./loader/test_folder/test_valid_program
	go test ./encoder -v -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html
	