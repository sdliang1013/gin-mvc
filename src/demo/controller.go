package demo

import (
	"errors"
	"gin-mvc/src/caul"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Controller struct {
}

// Routes 自定义路由信息
func (ctrl *Controller) Routes() []caul.Route {
	return []caul.Route{
		{Path: "/", Func: ctrl.GetData, Method: http.MethodGet},
		{Path: "/create", Func: ctrl.PostData, Method: http.MethodPost},
		{Path: "/download", Func: ctrl.GetDownload, Method: http.MethodGet},
	}
}

func (ctrl *Controller) GetData(ctx *gin.Context) (data any, err error) {
	// get param
	var id string
	// id 非空校验
	id, err = caul.NotNilString(ctx, "id")
	data = map[string]string{
		"id": id,
	}
	return
}

func (ctrl *Controller) PostData(ctx *gin.Context) (data any, err error) {
	// get body
	var body struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	}
	// 绑定body
	err = ctx.BindJSON(&body)
	data = body
	return
}

// GetDownload 文件下载
func (ctrl *Controller) GetDownload(ctx *gin.Context) (data any, err error) {
	// get param
	var path string
	path, err = caul.NotNilString(ctx, "path")
	if err != nil {
		return
	}
	// check file
	var stat os.FileInfo
	stat, err = os.Stat(path)
	if err != nil {
		return
	}
	if stat.IsDir() {
		err = errors.New("目标是文件夹")
		return

	}
	// download
	ctx.FileAttachment(path, stat.Name())
	return
}
