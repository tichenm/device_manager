package service

import (
	"fmt"
	"testing"
	//"zhiyuan/device_server/raying_api/internal/service"
)
func TestFeatureGateOverride(t *testing.T) {

	//hc := service.Httpclient{}
	username := "test"
	password := "qwerty123"
	method := "GET"
	url :=	"/api/cgi-bin/subscribe/picture"
	realm := "Login to 8L94R080029"
	nonce := "1753428206"
	nc :="00000001"
	cnonce := "ksjdfljwofsldj4687skjd"
	qop :="auth"
	HA1 := cal_md5(username+":"+realm+":"+password)
	HA2 := cal_md5(method+":"+url)
	response := cal_md5(HA1+":"+nonce+":"+nc+":"+cnonce+":"+qop+":"+HA2)
	fmt.Println(response)
	fmt.Println(HA2)
	fmt.Println(HA1)
}
