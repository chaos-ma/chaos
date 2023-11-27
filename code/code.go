package code

import (
	"net/http"

	"github.com/novalagung/gubrak"

	"github.com/chaos-ma/chaos/errors"
)

type ErrCode struct {
	C    int    //错误码
	HTTP int    //http的状态码
	Ext  string //扩展字段
	Ref  string //引用文档
}

func (e ErrCode) HTTPStatus() int {
	return e.HTTP
}

func (e ErrCode) String() string {
	return e.Ext
}

func (e ErrCode) Reference() string {
	return e.Ref
}

func (e ErrCode) Code() int {
	if e.C == 0 {
		return http.StatusInternalServerError
	}
	return e.C
}

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}
	var ref string
	if len(refs) > 0 {
		ref = refs[0]
	}
	coder := ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  ref,
	}

	errors.MustRegister(coder)
}

var _ errors.Coder = (*ErrCode)(nil)
