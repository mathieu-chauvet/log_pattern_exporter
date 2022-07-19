
VERSION=1.0.1-RELEASE

hello:
	echo "Hello"

build:
	#env GOOS=windows GOARCH=amd64 go build -o "bin/scp_copy_$(VERSION).exe" src/app/Server.go
	go test ./...
	env GOOS=linux GOARCH=amd64 go build -o "bin/log_pattern_exporter_$(VERSION)" src/app/exporter.go

run:
	go run main.go
