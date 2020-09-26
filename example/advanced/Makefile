
default: run

run:
	docker-compose up

proto:
	protoc -I$(GOPATH)/src --proto_path=. --go_out=${GOPATH}/src ./event/payload.proto
	./gorm.sh ./event/payload.pb.go
