package http

import (
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhiyuan/scaffold/internal/model"
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


	if err != nil || camera_params.Camera_address == "" || camera_params.Camera_position ==""||camera_params.Camera_type=="" {
		//if c.BindJSON(&json_subject)==nil{
		//	subject.Name = json_subject.Name
		//	subject.Photo = json_subject.Photo
		//	subject.Subject_Id = json_subject.Subject_Id
		//	subject.Subject_type = json_subject.Subject_type
		//	subject.Ip = json_subject.Ip
		//}else{
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		log4go.Error(err.Error())
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
		log4go.Error(err.Error())
	}

	if err != nil || camera_params.Camera_address == "" || camera_params.Camera_position ==""||camera_params.Camera_type=="" {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		log4go.Error(err.Error())

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
		log4go.Error(err.Error())
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

	if camera_params.Page == 0 {
		camera_params.Page = 1
	}
	if camera_params.Size == 0 {
		camera_params.Size = 10
	}
	if err != nil  {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		log4go.Error(err.Error())
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	res,err,count,total := svc.GetCameras(camera_params.Camera_status,camera_params.Page,camera_params.Size)

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
