# Chaos云原生微服务框架

## 1. grpc

1. 安装protobuf：https://github.com/protocolbuffers/protobuf/releases
2. 安装protoc-gen-go

```shell
go install github.com/golang/protobuf/protoc-gen-go@latest
```

3. 安装protoc-gen-go-grpc

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## 2. 模块

### 2.1 app

### 2.2 registry

### 2.3 server

#### 2.3.1 rpcserver

```
Interceptor借鉴go-zero的实现
```

#### 2.3.2  httpserver

### 2.4 log

```
基于zap封装的通用log接口
```

### 2.5 metadata

```
借鉴go-kratos的metadata实现
```

### 2.6 utils

### 2.7 third_party