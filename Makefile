build:
	CGO_ENABLED=0 go build -o ./mssql_exporter ./cmd

test:
	go test -v ./...