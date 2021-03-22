build:
	go build -o bin/compactor cmd/main.go

deploy: build
	sudo cp bin/compactor /usr/local/bin/compactor
	sudo chmod +x /usr/local/bin/compactor

run:
	go run cmd/main.go