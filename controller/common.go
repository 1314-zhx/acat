package controller

import (
	"acat/serializer"
	"encoding/json"
)

// ErrorResponse 封装错误返回，先判断是不是标志数据库的Json格式不正确错误。如果不是默认为参数异常
// 不对外暴露业务码，而是统一返回400
func ErrorResponse(err error) serializer.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: 400,
		Msg:    "参数错误",
		Error:  err.Error(),
	}
}
