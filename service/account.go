package service

import (
	"errors"
	"zhiyuan/scaffold/internal/model"
)

func (s *Service) UpdateAccount(AccountParams model.Account_josn,id int)(result model.Account,err error){

	//数据交换
	//传入DB方法

	err = s.dao.CheckAccount()
	if err != nil{
		return model.Account{},errors.New("初始化失败")
	}
	Add_Account := model.Account{
		ID:id,
		Account:AccountParams.Account,
		Password:AccountParams.Password,
		Ip_address:AccountParams.Ip_address,
	}
	obj,err := s.dao.UpdateAccount(Add_Account,Add_Account.ID)

	if err!=nil{
		return	model.Account{},err
	}
	return obj,nil
}

func (s *Service) GetAccount()(result model.Account,err error){
	err = s.dao.CheckAccount()
	if err != nil{
		return model.Account{},errors.New("初始化失败")
	}
	result,err = s.dao.GetAccount()

	if err !=nil{
		return model.Account{},err
	}
	return result,nil
}