install: 
	go install ./...


lint:
	golangci-lint run
