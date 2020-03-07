package dao

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"os"
	"zhiyuan/scaffold/internal/model"
)


func(d *Dao) CheckAccount()(error){
	var count int
	DBcount := d.crmdb
	if err := DBcount.Model(model.Account{}).Count(&count); err.Error != nil {
		return err.Error
	}
	if count == 0{
		data := model.Account{
			Account:"",
			Ip_address:"",
			Password:"",
		}
		if err := d.crmdb.Create(&data);err.Error!=nil{
			log.WithFields(log.Fields{
				"Account": "insert",
			}).Error("Account insert db err")
			return  err.Error
		}
		return nil
	}
	return  nil
}

func (d *Dao) UpdateAccount(data model.Account,id int)(Account_obj model.Account,err error){

	if err := d.crmdb.Model(&model.Account{}).Where("id = ?",id).Updates(data);err.Error!=nil{
		log.WithFields(log.Fields{
			"Account": "update",
		}).Error("Account update db err")
		return model.Account{}, err.Error
	}

		file_path :="./koala.toml"
		f, err := os.Create(file_path)
		if err != nil {
			// failed to create/open the file
			log.Fatal(err)
		}
		if err := toml.NewEncoder(f).Encode(data); err != nil {
			// failed to encode
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			// failed to close the file
			log.Fatal(err)
		}
	if err := d.crmdb.Model(&model.Account{}).Where("id = ?",id).Last(&Account_obj);err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "select",
		}).Error("select updated account in last time err")
		return model.Account{}, err.Error
	}
	return Account_obj,nil
}

func (d *Dao) GetAccount()(Account_obj model.Account,err error){

	if err := d.crmdb.Model(&model.Account{}).First(&Account_obj);err.Error!=nil{
		log.WithFields(log.Fields{
			"Camera": "select",
		}).Error("select updated account in last time err")
		return model.Account{}, err.Error
	}
	return Account_obj,nil
}


