# 生成error code工具

* 生成gin模板

```shell
cd chaos/cmd/chaos-codegen
go build -o chaos-codegen.exe
```

* 将chaos-codegen.exe拷贝至go_path的bin目录下

*进入code目录执行
```shell
go generate
```

* //chaos-codegen -type=int
* //chaos-codegen -type=int -doc -output ./error_code_generated.md