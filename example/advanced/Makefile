
proto:
	protoc -I$(GOPATH)/src --proto_path=. --go_out=${GOPATH}/src ./event/payload.proto
	protoc -I$(GOPATH)/src --proto_path=. --gorm_out=${GOPATH}/src ./event/payload.proto
