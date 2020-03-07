package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"net"
	"net/http"
	"strconv"
	"zhiyuan/scaffold/internal/koala"
	"zhiyuan/scaffold/internal/model"
	"zhiyuan/scaffold/service"
	"zhiyuan/zyutil_v1.5"
)

func  CreateCamera(c *gin.Context) () {

	var (
		resp4Device model.Resp4Device
		camera_params     model.Camera_json
		data        interface{}
	)

	//arg := new(model.Person)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	//err := c.Bind(&arg)
	err := c.BindJSON(&camera_params)
	//}else{
	//if err!=nil{
	//	ip := c.PostForm("ip")
	//	photo := c.PostForm("photo")
	//	name := c.PostForm("name")
	//	subject_id := c.PostForm("subject_id")
	//	subject_type := c.PostForm("subject_type")
	//	subject_id_int, err := strconv.Atoi(subject_id)
	//	subject_type_int, err := strconv.Atoi(subject_type)
	//
	//	subject.Ip = ip
	//	subject.Photo = photo
	//	subject.Name = name
	//	subject.Subject_type = subject_type_int
	//	subject.Subject_Id = subject_id_int
	//
	//}


	if err != nil || camera_params.Camera_address == "" || camera_params.Camera_position ==""||camera_params.Camera_type==""||camera_params.Camera_RTSP == "" {
		//if c.BindJSON(&json_subject)==nil{
		//	subject.Name = json_subject.Name
		//	subject.Photo = json_subject.Photo
		//	subject.Subject_Id = json_subject.Subject_Id
		//	subject.Subject_type = json_subject.Subject_type
		//	subject.Ip = json_subject.Ip
		//}else{
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		//log4go.Error(err.Error())
		//}
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	res,err := svc.CreateCamera(camera_params)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = res
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = res

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
	})
}

func  UpdateCamera(c *gin.Context) () {
	var (
		resp4Device model.Resp4Device
		camera_params     model.Camera_json
		data        interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&camera_params)

	tmp := c.Param("id")
	id, err := strconv.Atoi(tmp)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		//log4go.Error(err.Error())
	}

	if err != nil || camera_params.Camera_address == "" || camera_params.Camera_position ==""||camera_params.Camera_type==""||camera_params.Camera_RTSP == "" {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		//log4go.Error(err.Error())

	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	res,err := svc.UpdateCamera(camera_params,id)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = res
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = res

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
	})
}

func  DeleteCamera(c *gin.Context) () {
	var (
		resp4Device model.Resp4Device
		data        interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	tmp := c.Param("id")
	id, err := strconv.Atoi(tmp)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		//log4go.Error(err.Error())
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	err = svc.DeleteCamera(id)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = ""
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
	})
}

func  GetCameras(c *gin.Context) () {
	var (
		resp4Device model.Resp4Device
		camera_params     model.Camera_status_json
		data        interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&camera_params)


	if err != nil  {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		//log4go.Error(err.Error())
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	res,err,count,total := svc.GetCameras(camera_params.Camera_position,camera_params.Camera_status,camera_params.Page,camera_params.Size)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = res
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = res
	resp4Device.Page = map[string]interface{}{
		"count":count,
		"current":camera_params.Page,
		"total":total,
		"size":camera_params.Size,
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
		"page":resp4Device.Page,
	})
}

func  UpdateAccount(c *gin.Context) () {
	var (
		resp4Device model.Resp4Device
		account_params     model.Account_josn
		data        interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&account_params)

	tmp := c.Param("id")
	id, err := strconv.Atoi(tmp)
	//if err != nil {
	//	resp4Device.Code = -100
	//	resp4Device.Err_msg = "传入的参数有误!"
	//	//log4go.Error(err.Error())
	//}

	if err != nil || account_params.Account == "" || account_params.Password ==""||account_params.Ip_address=="" {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		//log4go.Error(err.Error())
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	res,err := svc.UpdateAccount(account_params,id)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = res
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = res

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
	})
}

func  GetAccount(c *gin.Context) () {
	var (
		resp4Device model.Resp4Device
		data        interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	res,err :=svc.GetAccount()

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = res
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = res
	resp4Device.Page = map[string]interface{}{}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
		"page":resp4Device.Page,
	})
}

func Get_display_config (c *gin.Context)(){

	var (
		resp4Device model.Resp4Device
		screenobj model.Screens
		screens_arr []model.Screens
	)

	device_token := c.Query("device_token")
	service.Device_token = device_token
	if device_token ==""{
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	rep,err := koala.GetDisplayDevice(device_token)
	if err!=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "与上位机失联"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	data :=rep.Get("data")
	device := data.Get("device")
	warning := data.Get("warning")
	yellowlist_warn := data.Get("yellowlist_warn")
	screens,_ := data.Get("screens").Array()
	res,err :=svc.GetAllCameras()
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询数据库失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	ip,_:=LocalIPv4s()
	if len(screens)!=0{
		mapstructure.Decode(screens[0].(map[string]interface{}), &screenobj)
		for i:= 0;i<len(res);i++{
			screenrecive := model.Screens{}
			screenrecive = screenobj
			screenrecive.Box_address = ip[0]
			screenrecive.Camera_address =	res[i].Camera_RTSP
			screenrecive.Camera_position = res[i].Camera_position
			screenrecive.Screen_token = res[i].Camera_token
			screenrecive.Camera_status = strconv.Itoa(res[i].Camera_status)
			screenrecive.Id = res[i].ID
			screenrecive.Is_select = 1
			screenrecive.Type = 1
			screens_arr = append(screens_arr, screenrecive)
		}
	}
	result := map[string]interface{}{
		"device":device,
		"screens":screens_arr,
		"warning":warning,
		"yellowlist_warn":yellowlist_warn,
	}
	resp4Device.Page = map[string]interface{}{}
	fmt.Println(result)
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"data":    result,
		"page":resp4Device.Page,
	})
}


func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips, nil
}

func Set_display_config (c *gin.Context)(){
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data": map[string]interface{}{} ,
		"page":map[string]interface{}{},
	})
}

func Get_screen_list(c *gin.Context)(){

	var (
		resp4Device model.Resp4Device
		screens_arr []model.Screens
	)

	box_token := c.Query("box_token")
	service.Box_token = box_token

	res,err :=svc.GetAllCameras()
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询数据库失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	//if len(screens)!=0{
		//mapstructure.Decode(screens[0].(map[string]interface{}), &screenobj)
	for i:= 0;i<len(res);i++{
		screenrecive := model.Screens{}
		//screenrecive = screenobj
		screenrecive.Camera_address =	res[i].Camera_RTSP
		screenrecive.Camera_position = res[i].Camera_position
		screenrecive.Screen_token = res[i].Camera_token
		screenrecive.Camera_status = strconv.Itoa(res[i].Camera_status)
		screenrecive.Id = res[i].ID
		screenrecive.Is_select = 1
		screenrecive.Type = 1
		screens_arr = append(screens_arr, screenrecive)
	}
	//}
	result := map[string]interface{}{
		"screens":screens_arr,
	}
	resp4Device.Page = map[string]interface{}{}
	fmt.Println(result)
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"data":    result,
		"page":resp4Device.Page,
	})
}