# protobuf 生成gin模板

* 生成gin模板

```shell
cd chaos/template
go build -o protoc-gen-gin.exe
```

* 将protoc-gen-gin.exe拷贝至go_path的bin目录下

* protoc命令

```shell
protoc --proto_path=. --proto_path=../third_party --go_out=./test --go-grpc_out=./test --gin_out=./test test.proto
```

