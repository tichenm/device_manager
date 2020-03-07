package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"zhiyuan/scaffold/internal/koala"
)

var key1 = []byte{'0', 0x01, 0x22}


func ReverseProxy(c *gin.Context) {
	//LogRequest(c)  // POST 不能写日志，不然代理会少内容
	koala.ReverseProxy2(c)
}
func ReverseProxy_koalamate(c *gin.Context) {
	//LogRequest(c)  // POST 不能写日志，不然代理会少内容
	koala.ReverseProxy3(c)
}

func KoalaStatic(c *gin.Context) {
	// 去掉url前面的 /koala_static
	c.Request.URL.Path = c.Param("path")
	koala.ReverseProxy2(c)

}

func GetServerTime(c *gin.Context) {

	var serverTime = time.Now().Unix()

	c.Header("Connection", "close")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"err_msg": "",
		"data": serverTime,
	})
}