test:
	go test -v -count=1 ./... -cover

lint:	
	golangci-lint run



