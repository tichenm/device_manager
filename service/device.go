package service

import (
	"errors"
	"github.com/satori/go.uuid"
	"zhiyuan/scaffold/internal/model"
)

func (s *Service) CreateCamera(CameraParams model.Camera_json)(result model.Camera,err error){

	//数据交换
	//传入DB方法

	uid, _ := uuid.NewV4()

	Add_Camera := model.Camera{
		Camera_type:CameraParams.Camera_type,
		Camera_position:CameraParams.Camera_position,
		Camera_address:CameraParams.Camera_address,
		Camera_RTSP:CameraParams.Camera_RTSP,
		Camera_status:2,
		Camera_token:uid.String(),
	}
	res,err := s.dao.CheckCameras(Add_Camera.Camera_address)
	if res != true{
		return model.Camera{},errors.New("ip重复")
	}
	obj,err := s.dao.AddCamera(Add_Camera)
	if err!=nil{
		return	model.Camera{},err
	}
	return obj,nil
}

func (s *Service) UpdateCamera(CameraParams model.Camera_json,id int)(result model.Camera,err error){

	//数据交换
	//传入DB方法
	Add_Camera := model.Camera{
		ID:id,
		Camera_type:CameraParams.Camera_type,
		Camera_position:CameraParams.Camera_position,
		Camera_address:CameraParams.Camera_address,
		Camera_RTSP:CameraParams.Camera_RTSP,
	}

	obj,err := s.dao.UpdateCamera(Add_Camera,Add_Camera.ID)
	if err!=nil{
		return	model.Camera{},err
	}
	return obj,nil
}

func (s *Service) DeleteCamera(id int)(err error){

	err = s.dao.DeleteCamera(id)
	if err!=nil{
		return	err
	}
	return nil
}
func (s *Service) GetCameras(camera_position string,camera_status,page,size int)(result []model.Camera,err error,count int,total int){

	result,err,count,total = s.dao.GetCameras(camera_position,camera_status,page,size)
	if err!=nil{
		return	[]model.Camera{},err,0,0
	}
	return result,nil,count,total
}

func (s *Service) GetAllCameras()(result []model.Camera,err error){

	result,err = s.dao.GetAllCameras()
	if err!=nil{
		return	[]model.Camera{},err
	}
	return result,nil
}






