.SILENT:

build-server:
	go build -o ./bin/server cmd/server/main.go

run-server:
	./bin/server

server:	build-server run-server


build-client:
	cd cmd/client && wails build && cp ./build/bin/rps-online ../../bin/client

run-client:
	./bin/client

client: build-client run-client