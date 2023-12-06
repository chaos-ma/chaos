# Chaos云原生微服务框架

### 1. 安装protobuf：https://github.com/protocolbuffers/protobuf/releases

### 2. 安装protoc-gen-go

```shell
go install github.com/golang/protobuf/protoc-gen-go@latest
```

### 3. 安装protoc-gen-go-grpc

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 4. 安装stringer

```shell
go install golang.org/x/tools/cmd/stringer@latest
```

### 5. protoc-gen-gin工具安装(protobuf 生成gin模板)

* 5.1 build生成工具

```shell
cd chaos/template
go build -o protoc-gen-gin.exe
```

* 5.2 将protoc-gen-gin.exe拷贝至go_path的bin目录下

* 5.3 protoc命令

```shell
protoc --proto_path=. --proto_path=../third_party --go_out=./test --go-grpc_out=./test --gin_out=./test test.proto
```

### 6. 生成error code工具

* 6.1 build生成工具

```shell
cd chaos/cmd/chaos-codegen
go build -o chaos-codegen.exe
```

* 6.2 chaos-codegen.exe拷贝至go_path的bin目录下
