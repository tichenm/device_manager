package dao

import (
	log "github.com/sirupsen/logrus"
	"zhiyuan/scaffold/internal/model"
	"zhiyuan/zyutil_v1.5"
)

func (d *Dao) AddCamera(data model.Camera)(Camera_obj model.Camera,err error){

	if err := d.crmdb.Create(&data);err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "insert",
		}).Error("camera insert db err")
		return model.Camera{}, err.Error
	}
	if err := d.crmdb.Last(&Camera_obj);err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "select",
		}).Error("select camera in last time err")
		return model.Camera{}, err.Error
	}
	return Camera_obj,nil
}


func (d *Dao) UpdateCamera(data model.Camera,id int)(Camera_obj model.Camera,err error){

	if err := d.crmdb.Model(&model.Camera{}).Where("id = ?",id).Updates(data);err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "update",
		}).Error("camera update db err")
		return model.Camera{}, err.Error
	}
	if err := d.crmdb.Model(&model.Camera{}).Where("id = ?",id).Last(&Camera_obj);err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "select",
		}).Error("select updated camera in last time err")
		return model.Camera{}, err.Error
	}
	return Camera_obj,nil
}

func (d *Dao) DeleteCamera(id int)( err error){
	if err := d.crmdb.Where("id = ?",id).Delete(&model.Camera{});err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "delete",
		}).Error("camera delete db err")
		return  err.Error
	}
	return nil
}

func(d *Dao) GetCameras(camera_pisition string,camera_status ,page , size int)(result []model.Camera,err error,count int,total int){

	DBdate := d.crmdb
	if camera_status != 0 {
		DBdate = DBdate.Where("camera_status = ?", camera_status)
	}
	if camera_pisition != "" {
		DBdate = DBdate.Where("camera_position = ?", camera_pisition)
	}
	DBdate = DBdate.Order("id desc")
	if page > 0 {
		DBdate = DBdate.Limit(size).Offset((page - 1) * size)
	}
	if err := DBdate.Model(model.Camera{}).Find(&result);err.Error!= nil {

		return result, err.Error, 0,0
	}
	DBcount := d.crmdb
	if camera_status != 0 {
		DBcount = DBcount.Where("camera_status = ?", camera_status)
	}
	if camera_pisition != "" {
		DBcount = DBcount.Where("camera_position = ?", camera_pisition)
	}
	if err := DBcount.Model(model.Camera{}).Count(&count); err.Error != nil {
		return result, err.Error, 0,0
	}
	total = zyutil.GetTotal(count, size)
	return result, nil, count, total
}

func(d *Dao) CheckCameras(camera_address  string)(bool,error){
	var count int
	DBcount := d.crmdb
	DBcount = DBcount.Where("camera_address = ?", camera_address)
	if err := DBcount.Model(model.Camera{}).Count(&count); err.Error != nil {
		return false,err.Error
	}
	if count != 0{
		return false,nil
	}
	return true, nil
}

func(d *Dao) GetAllCameras()(result []model.Camera,err error){

	DBdate := d.crmdb
	DBdate = DBdate.Order("id desc")
	if err := DBdate.Model(model.Camera{}).Find(&result);err.Error!= nil {
		return result, err.Error
	}
	return result, nil
}