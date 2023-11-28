type {{$.Name}}HttpServer struct{
	server {{ $.ServiceName }}
	router gin.IRouter
}

func Register{{ $.ServiceName }}HTTPServer(srv {{ $.ServiceName }}, r gin.IRouter) {
	s := {{.Name}}HttpServer{
		server: srv,
		router:     r,
	}
	s.RegisterService()
}

{{range .Methods}}
func (s *{{$.Name}}HttpServer) {{ .HandlerName }} (c *gin.Context) {
	var in {{.Request}}
{{if eq .Method "GET" "DELETE" }}
	if err := c.ShouldBindQuery(&in); err != nil {
		s.resp.ParamsError(ctx, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
{{else}}
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
{{end}}
{{if .HasPathParams }}
	{{range $item := .PathParams}}
	in.{{$.GoCamelCase $item }} = c.Params.ByName("{{$item}}")
	{{end}}
{{end}}
	out, err := s.server.{{.Name}}(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	c.JSON(http.StatusOK, out)
}
{{end}}

func (s *{{$.Name}}HttpServer) RegisterService() {
{{range .Methods}}
		s.router.Handle("{{.Method}}", "{{.Path}}", s.{{ .HandlerName }})
{{end}}
}