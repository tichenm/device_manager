package service

import (
	"github.com/Lofanmi/pinyin-golang/pinyin"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
	"time"
	"zhiyuan/device_server/raying_api/internal/model"
)

func koalaadperson(json1 *simplejson.Json, PersonInfo model.RYPerson_2, photo string, subject_type int) (map[string]interface{}, error) {
	var facetoken string
	_, ok := json1.Get("params").CheckGet("FaceToken")
	if ok {
		FaceToken_list, _ := json1.Get("params").Get("FaceToken").Array()
		facetoken = FaceToken_list[0].(string)
	} else {
		facetoken = ""
	}

	dict := pinyin.NewDict()
	pinyin := dict.Sentence(PersonInfo.Name).None()
	pinyin = strings.Replace(pinyin, " ", "", -1)
	sub_id, _ := strconv.Atoi(PersonInfo.ID)
	person_photos := model.Koala_person_photos{
		Url:        photo,
		Company_id: 1,
		Id:         sub_id,
		Subject_id: sub_id,
		Version:    1,
		Quality:    0.99453,
		Origin_url: facetoken,
	}
	photos := []model.Koala_person_photos{person_photos}
	timestamp := time.Now().Unix()
	koala_person_info := model.Koala_person{
		Avatar:             photo,
		Birthday:           timestamp,
		Come_from:          "",
		Company_id:         1,
		Create_time:        timestamp,
		Department:         "",
		Description:        "",
		Email:              "",
		End_time:           timestamp,
		Entry_date:         timestamp,
		Gender:             0,
		Id:                 sub_id,
		Interviewee:        "",
		Interviewee_pinyin: "",
		Job_number:         "",
		Name:               PersonInfo.Name,
		Password_reseted:   false,
		Phone:              "",
		Photos:             photos,
		Pinyin:             pinyin,
		Purpose:            0,
		Remark:             "",
		Start_time:         timestamp,
		Subject_type:       subject_type,
		Title:              "",
		Visit_notify:       false,
	}
	return_data := map[string]interface{}{
		"code":    0,
		"err_msg": "",
		"data":    koala_person_info,
	}
	return return_data, nil
}

func Getfacetoken(json1 *simplejson.Json) (*simplejson.Json) {
	var FaceToken_list *simplejson.Json
	_, ok := json1.Get("params").CheckGet("FaceToken")
	if ok {
		FaceToken_list = json1.Get("params")
	}

	return FaceToken_list
}

//func koalaadphoto(json1 *simplejson.Json,PersonInfo model.RYPerson_2,photo string,subject_type int)(model.Koala_person,error){
//	var facetoken string
//	_, ok := json1.Get("params").CheckGet("FaceToken")
//	if ok {
//		FaceToken_list,_ := json1.Get("params").Get("FaceToken").Array()
//		//fmt.Println(len(FaceToken_list))
//		facetoken = FaceToken_list[0].(string)
//		//fmt.Println(FaceToken_list[0])
//	}else{
//		facetoken = ""
//	}
//
//	dict := pinyin.NewDict()
//	pinyin := dict.Sentence(PersonInfo.Name).None()
//	pinyin = strings.Replace(pinyin, " ", "", -1)
//	sub_id, _ := strconv.Atoi(PersonInfo.ID)
//	person_photos:=model.Koala_person_photos{
//		Url:photo,
//		Company_id:1,
//		Id:sub_id,
//		Subject_id:sub_id,
//		Version:1,
//		Quality:0.99453,
//		Origin_url:facetoken,
//	}
//	photos := []model.Koala_person_photos{person_photos}
//	timestamp := time.Now().Unix()
//	koala_person_info := model.Koala_person{
//		Avatar:photo,
//		Birthday:timestamp,
//		Come_from:"",
//		Company_id:1,
//		Create_time:timestamp,
//		Department:"",
//		Description:"",
//		Email:"",
//		End_time:timestamp,
//		Entry_date:timestamp,
//		Gender:0,
//		Id:sub_id,
//		Interviewee:"",
//		Interviewee_pinyin:"",
//		Job_number:"",
//		Name:PersonInfo.Name,
//		Password_reseted:false,
//		Phone:"",
//		Photos:	photos,
//		Pinyin:pinyin,
//		Purpose:0,
//		Remark:"",
//		Start_time:timestamp,
//		Subject_type:subject_type,
//		Title:"",
//		Visit_notify:false,
//	}
//
//	return koala_person_info,nil
//}
