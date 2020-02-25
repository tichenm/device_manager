package http

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"runtime"
	"zhiyuan/scaffold/internal/model"
	"zhiyuan/scaffold/service"
	"zhiyuan/scaffold/configs"
)
var(
	svc *service.Service

)
//var (
//	svc *service.Service
//)

// New new a gin server.
func New() {
	var (
		hc struct {
			Server *model.ServerConfig
		}
		cpath string
		//screen struct{
		//	ac  *paladin.Map
		//}
	)

	// 初始化
	if runtime.GOOS == "linux" {
		cpath = "./configs/http.toml"
		if runtime.GOARCH == "arm" {
			cpath = "./configs/http.toml"
		} else if runtime.GOARCH == "amd64" {
			cpath = "./configs/http.toml"
		}
	} else if runtime.GOOS == "windows" {
		cpath = "F:/program/code/go/go_project/src/zhiyuan/device_server/raying_api/configs/http.toml"
	}

	if _, err := toml.DecodeFile(cpath, &hc); err != nil {
		log.Error("read toml file error(%v)", err)
	}

	svc = service.New(configs.Conf)
	engine := gin.Default()
	initRouter(engine)
	gin.SetMode(gin.ReleaseMode)
	engine.Run(hc.Server.Addr)
}

func initRouter(e *gin.Engine) {

	system := e.Group("/v1")
	{
		system.GET("/start", howToStart)
	}
	service := e.Group("/v1")
	{
		service.POST("/system/screen",CreateCamera)
		service.PUT("/system/screen/:id",UpdateCamera)
		service.DELETE("/system/screen/:id",DeleteCamera)
		service.GET("/system/screen",GetCameras)


	}
}

// example for http request handler.
func howToStart(c *gin.Context) {
	c.String(0, "Golang 大法好 !!!")
}
