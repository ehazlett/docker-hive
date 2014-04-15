all:
	@go build .
	@go get -d ./...

deps:
	@go get -d ./...

build:
	@go build .
benchmark:
	@go test ./hive -bench=.
test:
	@go test ./...
fmt:
	@go fmt .
	@go fmt ./hive/
clean:
	@rm -rf docker-hive
