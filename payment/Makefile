.SILENT:

proto:
	protoc --proto_path=./internal/proto-files/statistics stat.proto --go_out=plugins=grpc:./internal/proto-files/statistics

run_pay:
	go build -o=./bin/payment ./cmd/payment/main.go
	./bin/payment

run_stat:
	go build -o=./bin/stat ./cmd/statserver/main.go
	./bin/stat
