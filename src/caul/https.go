package caul

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	_httpErr HttpError
)

const (
	ErrCodeOk      = "0"
	ErrCodeUnknown = "-1"
)

type ErrCode string

type ResponseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HttpError interface {
	Status() int
	Error
}

type ResponseFunction func(ctx *gin.Context, data any, err HttpError)

func GetJson(ctx *gin.Context) (data map[string]interface{}, err error) {
	err = ctx.BindJSON(&data)
	return
}

func QueryString(ctx *gin.Context, name string) string {
	return ctx.Query(name)
}

func QueryBool(ctx *gin.Context, name string) (bool, error) {
	return strconv.ParseBool(ctx.DefaultQuery(name, "false"))
}

func QueryFloat(ctx *gin.Context, name string) (float64, error) {
	return strconv.ParseFloat(ctx.DefaultQuery(name, "0"), 64)
}

func QueryInt(ctx *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(ctx.DefaultQuery(name, "0"), 10, 64)
}

func QueryUint(ctx *gin.Context, name string) (uint64, error) {
	return strconv.ParseUint(ctx.DefaultQuery(name, "0"), 10, 64)
}

func NotNilString(ctx *gin.Context, name string) (param string, err error) {
	param = ctx.Query(name)
	if param == "" {
		err = fmt.Errorf("%s 不能为空", name)
		return
	}
	return
}

func NotNilBool(ctx *gin.Context, name string) (param bool, err error) {
	str := ctx.Query(name)
	if str == "" {
		err = fmt.Errorf("%s 不能为空", name)
		return
	}
	param, err = strconv.ParseBool(str)
	return
}

func NotNilFloat(ctx *gin.Context, name string) (param float64, err error) {
	str := ctx.Query(name)
	if str == "" {
		err = fmt.Errorf("%s 不能为空", name)
		return
	}
	param, err = strconv.ParseFloat(str, 64)
	return
}

func NotNilInt(ctx *gin.Context, name string) (param int64, err error) {
	str := ctx.Query(name)
	if str == "" {
		err = fmt.Errorf("%s 不能为空", name)
		return
	}
	param, err = strconv.ParseInt(str, 10, 64)
	return
}

func NotNilUint(ctx *gin.Context, name string) (param uint64, err error) {
	str := ctx.Query(name)
	if str == "" {
		err = fmt.Errorf("%s 不能为空", name)
		return
	}
	param, err = strconv.ParseUint(str, 10, 64)
	return
}

func ResponseString(ctx *gin.Context, data string, err HttpError) {
	// check err
	if err != nil {
		ctx.String(err.Status(), err.Error())
		return
	}
	//response
	ctx.String(http.StatusOK, data)
}

func ResponseJson(ctx *gin.Context, data any, err HttpError) {
	// check err
	if err != nil {
		ctx.JSON(err.Status(), ErrBody(err.Code(), err.Error(), nil))
		return
	}
	//response
	ctx.JSON(http.StatusOK, OkBody(data))
}

func ResponseXml(ctx *gin.Context, data any, err HttpError) {
	// check err
	if err != nil {
		ctx.XML(err.Status(), ErrBody(err.Code(), err.Error(), nil))
		return
	}
	//response
	ctx.XML(http.StatusOK, OkBody(data))
}

func ResponseYaml(ctx *gin.Context, data any, err HttpError) {
	// check err
	if err != nil {
		ctx.YAML(err.Status(), ErrBody(err.Code(), err.Error(), nil))
		return
	}
	//response
	ctx.YAML(http.StatusOK, OkBody(data))
}

func OkBody(data any) (body ResponseBody) {
	return ResponseBody{
		Code:    ErrCodeOk,
		Message: "success",
		Data:    data,
	}
}

func ErrBody(code string, message string, data any) (body ResponseBody) {
	return ResponseBody{
		Code:    DefaultString(code, ErrCodeUnknown),
		Message: DefaultString(message, "error"),
		Data:    data,
	}
}

type HttpErrorWrapper struct {
	status int
	err    Error
}

func (e HttpErrorWrapper) Status() int {
	return e.status
}

func (e HttpErrorWrapper) Cause() error {
	return e.err.Cause()
}

func (e HttpErrorWrapper) Code() string {
	return e.err.Code()
}

func (e HttpErrorWrapper) Message() string {
	return e.err.Message()
}

func (e HttpErrorWrapper) Error() string {
	return e.err.Error()
}

func WrapError(err error) HttpError {
	if err == nil {
		return nil
	}
	// HttpError
	if errors.As(err, &_httpErr) {
		return err.(HttpError)
	}
	// Error
	if errors.As(err, &_err) {
		return HttpErrorWrapper{
			status: http.StatusInternalServerError,
			err:    err.(Error),
		}
	}
	// 包装
	return HttpErrorWrapper{
		status: http.StatusInternalServerError,
		err:    NewError(ErrCodeUnknown, err.Error(), nil),
	}
}

func NewHttpError(status int, code string, message string, err error) HttpError {
	// Error
	if errors.As(err, &_err) {
		return HttpErrorWrapper{
			status: DefaultInt(status, http.StatusInternalServerError),
			err:    err.(Error),
		}
	}
	// 包装
	return HttpErrorWrapper{
		status: DefaultInt(status, http.StatusInternalServerError),
		err:    NewError(DefaultString(code, ErrCodeUnknown), message, err),
	}
}

func Goroutine(ctx *gin.Context, goroutine func(ctx *gin.Context)) {
	// 创建在 goroutine 中使用的副本
	goroutine(ctx.Copy())
}
