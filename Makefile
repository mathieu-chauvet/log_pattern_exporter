
#VERSION=1.0.3-SNAPSHOT
VERSION=1.0.3-RELEASE

hello:
	echo "Hello"

build:
	#env GOOS=windows GOARCH=amd64 go build -o "bin/scp_copy_$(VERSION).exe" src/app/Server.go
	go test ./...
	env GOOS=linux GOARCH=amd64 go build -o "bin/log_pattern_exporter_$(VERSION)" src/app/exporter.go
	cp "bin/log_pattern_exporter_$(VERSION)" "bin/log_pattern_exporter"
	tar -czvf "bin/log_pattern_exporter-$(VERSION).linux-amd64.tar.gz"  -C bin "log_pattern_exporter"
run:
	go run main.go
