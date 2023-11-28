package main

import (
	"flag"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/chaos-ma/chaos/cmd/protoc-gen-gin/generator"
)

//go:generate protoc --proto_path=. --proto_path=../third_party --go_out=./test --go-grpc_out=./test --gin_out=./test test.proto
func main() {
	flag.Parse()
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generator.GenerateFile(gen, f)
		}
		return nil
	})
}
