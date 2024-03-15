package caul

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"reflect"
	"strings"
)

const (
	PathSplit   = "/"
	PkgSplit    = "."
	ContentType = "json"
)

var (
	methods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
)

// MidHandler 中间件
type MidHandler = gin.HandlerFunc

// RouteFunc 路由方法签名
// todo: 支持任意形式的method
type RouteFunc func(ctx *gin.Context) (data any, err error)

// Route 路由信息
type Route struct {
	Path   string
	Method string
	Func   RouteFunc
}

// Controller 控制器
type Controller interface {
	Routes() []Route //路由信息
}

// CRoute 控制器路由
type CRoute struct {
	Path       string       // 上级路由
	Controller Controller   //注册的接口
	Handlers   []MidHandler //中间件
}

type CRouter struct {
	gin.IRouter
}

// RegisterMiddleware 注册中间件
func (r *CRouter) RegisterMiddleware(handlers ...MidHandler) CRouter {
	r.Use(handlers...)
	return *r
}

// RegisterRoute 注册控制器
//
//	@param cRoute 控制器路由
func (r *CRouter) RegisterRoute(cRoute CRoute) CRouter {
	if !strings.HasPrefix(cRoute.Path, PathSplit) {
		panic(fmt.Errorf("Path must start with %s: %s\n", PathSplit, cRoute.Path))
	}
	group := r.Group(cRoute.Path, cRoute.Handlers...)
	// scan method
	for _, route := range cRoute.Controller.Routes() {
		if !strings.HasPrefix(route.Path, PathSplit) {
			panic(fmt.Errorf("Path must start with %s: %s\n", PathSplit, route.Path))
		}
		group.Handle(route.Method, route.Path, r.FuncWrapper(route.Func))
	}
	return *r
}

// RegisterController 注册控制器
//
//	@param [relativePath 上级路径, controller 控制器, handlers 插件]
func (r *CRouter) RegisterController(relativePath string, controller interface{}, handlers ...MidHandler) CRouter {
	if !strings.HasPrefix(relativePath, PathSplit) {
		panic(fmt.Errorf("relativePath must start with %s: %s\n", PathSplit, relativePath))
	}
	group := r.Group(relativePath, handlers...)
	ctrlStruct := reflect.ValueOf(controller)
	ctrlType := reflect.TypeOf(controller)
	// include package path
	//relativePath = r.wrapPkgPath(relativePath, ctrlType)
	// scan method
	for i := 0; i < ctrlStruct.NumMethod(); i++ {
		funcName := ctrlType.Method(i).Name
		httpMethod, valid := r.getHttpMethod(funcName)
		// not valid method
		if !valid {
			continue
		}
		method := ctrlStruct.Method(i)
		realPath := r.getPath(funcName)
		handler := func(method reflect.Value) gin.HandlerFunc {
			return func(c *gin.Context) {
				r.parseResults(c, method.Call(r.parseParams(c, method)))
			}
		}(method)
		group.Handle(httpMethod, realPath, handler)
	}
	return *r
}

// FuncWrapper 原始方法包装
func (r *CRouter) FuncWrapper(method RouteFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// todo: 支持任意形式的method
		res, err := method(ctx)
		if err != nil {
			r.responseFunc()(ctx, res, WrapError(err))
			return
		}
		if res != nil {
			r.responseFunc()(ctx, res, nil)
		}
	}
}

// wrapPkgPath 包路径
func (r *CRouter) wrapPkgPath(relativePath string, ctrlType reflect.Type) string {
	ctrlPath := ctrlType.String()
	ctrlName := ctrlType.Name()
	if ctrlPath == ctrlName {
		return relativePath
	}
	ctrlPath = ctrlPath[0 : len(ctrlPath)-len(ctrlName)-1]
	return path.Join(relativePath, strings.ReplaceAll(ctrlPath, PkgSplit, PathSplit))
}

// getHttpMethod 获取HTTP方法
//
//	@param funcName 以HttpMethod开头 GetXXX PostXXX...
func (r *CRouter) getHttpMethod(funcName string) (string, bool) {
	words := SplitCameCase(funcName)
	method := strings.ToUpper(words[0])
	if ContainsString(methods, method) {
		return method, true
	}
	return "", false
}

// getPath 获取请求路径
//
//	驼峰命名法 GetAbcDef -> GET请求, 路径/abc/def
//	@param [relativePath 上级路径, funcName 以HttpMethod开头 GetAbcDef PostAbcDef]
func (r *CRouter) getPath(funcName string) string {
	words := SplitCameCase(funcName)
	if len(words) == 1 {
		return ""
	}
	return strings.ToLower(path.Join(words[1:]...))
}

// parseParams 参数处理
func (r *CRouter) parseParams(ctx *gin.Context, method reflect.Value) []reflect.Value {
	// todo: 反射 method的参数列表 的值
	return []reflect.Value{reflect.ValueOf(ctx)}
}

// parseResults 结果处理
func (r *CRouter) parseResults(ctx *gin.Context, results []reflect.Value) {
	switch len(results) {
	case 0:
		return
	case 1:
		r.responseFunc()(ctx, r.parseData(results[0]), nil)
	default:
		r.responseFunc()(ctx, r.parseData(results[0]), r.parseError(results[1]))
	}
}

// parseData 处理业务数据
func (r *CRouter) parseData(result reflect.Value) any {
	if result.CanInt() {
		return result.Int()
	}
	if result.CanFloat() {
		return result.Float()
	}
	if result.CanUint() {
		return result.Uint()
	}
	if result.CanAddr() {
		return result.Addr()
	}
	if result.CanComplex() {
		return result.Complex()
	}
	if result.CanInterface() {
		return result.Interface()
	}
	return result.String()
}

// parseData 处理异常
func (r *CRouter) parseError(result reflect.Value) HttpError {
	if result.IsNil() {
		return nil
	}
	err, isErr := result.Interface().(error)
	if !isErr || err == nil {
		return nil
	}
	return WrapError(err)
}

// responseFunc
func (r *CRouter) responseFunc() ResponseFunction {
	switch ContentType {
	case "xml":
		return ResponseXml
	case "yaml":
		return ResponseYaml
	default:
		return ResponseJson
	}
}
