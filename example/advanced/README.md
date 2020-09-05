## Advanced Tributary Example

Install `protoc-gen-gorm`

```
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/infobloxopen/protoc-gen-gorm
cd $GOPATH/src/github.com/infobloxopen/protoc-gen-gorm
go install
```

In this directory, generate the proto files if changed

```
make proto
```
