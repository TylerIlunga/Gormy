build:
	rm -rf ./bin ./vendor Gopkg.lock
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/gormy ./main.go
.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock
