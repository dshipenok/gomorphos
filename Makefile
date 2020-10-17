
test:
	go test -cover ./...

coverage:
	go test -coverprofile=cov.out ./...
	go tool cover -html=cov.out
