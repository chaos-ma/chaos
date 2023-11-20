package generator

/**
* created by mengqi on 2023/11/15
 */

import (
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	ginPkg  = protogen.GoImportPath("github.com/gin-gonic/gin")
	httpPkg = protogen.GoImportPath("net/http")
)

var methodSets = make(map[string]int)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	//设置生成的文件名，文件名会被protoc使用，生成的文件会被放在响应的目录下
	filename := file.GeneratedFilenamePrefix + "_gin.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	//该注释会被go的ide识别到， 表示该文件是自动生成的，尽量不要修改
	g.P("// Code generated by protoc-gen-gin. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)

	//该函数是注册全局的packge 的内容，但是此时不会写入
	g.QualifiedGoIdent(ginPkg.Ident(""))
	g.QualifiedGoIdent(httpPkg.Ident(""))

	for _, service := range file.Services {
		genService(file, g, service)
	}

	//自己写文件看结果
	//f, err := os.Create("api_gin.pb.go")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer f.Close()
	//
	//contentStr, _ := g.Content()
	//_, _ = f.WriteString(string(contentStr))

	return g
}

func genService(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	// HTTP Server
	sd := &service{
		Name:     s.GoName,
		FullName: string(s.Desc.FullName()),
	}

	for _, method := range s.Methods {
		sd.Methods = append(sd.Methods, genMethod(method)...)
	}

	text := sd.execute()
	g.P(text)
}

func genMethod(m *protogen.Method) []*method {
	var methods []*method

	// 存在 http rule 配置
	// options
	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule != nil && ok {
		//for _, bind := range rule.AdditionalBindings {
		//	methods = append(methods, buildHTTPRule(m, bind))
		//}

		methods = append(methods, buildHTTPRule(m, rule))
		return methods
	}

	methods = append(methods, defaultMethod(m))
	return methods
}

func defaultMethod(m *protogen.Method) *method {
	// TODO path
	// $prefix + / + ${package}.${service} + / + ${method}
	// /api/demo.v0.Demo/GetName
	md := buildMethodDesc(m, "POST", "")
	md.Body = "*"
	return md
}

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *method {
	var path, method string
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}

	md := buildMethodDesc(m, method, path)
	return md
}

func buildMethodDesc(m *protogen.Method, httpMethod string, path string) *method {
	defer func() { methodSets[m.GoName]++ }()

	md := &method{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: m.Input.GoIdent.GoName,
		Reply:   m.Output.GoIdent.GoName,
		Path:    path,
		Method:  httpMethod,
	}

	md.initPathParams()
	return md
}
