build:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/go_api main.go 

run: build
	@./bin/go_api