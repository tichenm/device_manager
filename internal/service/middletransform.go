package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"zhiyuan/device_server/raying_api/internal/model"
)

func digital2String(digital string) string {
	switch (digital) {
	case "589825":
		return "未知错误"
	case "589826":
		return "人脸底库组名不唯一"
	case "589827":
		return "人脸底库组别名不唯一"
	case "589828":
		return "人脸底库组GUID不唯一"
	case "589829":
		return "人脸底库组数已达上限"
	case "589830":
		return "入库人员名字不唯一"
	case "589831":
		return "入库人员证件和证件号不唯一"
	case "589832":
		return "人脸图片Token不唯一"
	case "589833":
		return "人脸图片重复绑定"
	case "589834":
		return "入库人脸数已达上限"
	case "589835":
		return "此人入库人脸数已达上限"
	case "589836":
		return "此人待入库人脸数超过上限"
	case "589837":
		return "人脸图片数据空"
	case "589838":
		return "人员不存在"
	case "589839":
		return "人脸不存在"
	case "589840":
		return "人脸图片已入库"
	case "589841":
		return "人脸图片无人脸"
	case "589842":
		return "人脸图片质量过低"
	case "589843":
		return "人脸底库组不存在"
	case "589844":
		return "人脸图片数据过大"
	case "589845":
		return "系统忙"
	case "589846":
		return "算法未知错误"
	case "589847":
		return "模型版本不匹配"
	default:
		return "UNKNOWN"
	}
}

func Init(ip, session, authorization string) Rayingclient {
	hc := CreateConnection(ip, session, authorization)
	log.WithFields(log.Fields{
		"raying-cron": hc.Host,
	}).Info("start Checkauth")
	flag := hc.Checkauth()
	if flag {
		log.WithFields(log.Fields{
			"raying-cron": hc.Host,
		}).Info("Checkauth success first time")
		return hc
	} else {
		flag := Retry(10000, 30, hc.Checkauth)
		log.WithFields(log.Fields{
			"raying-cron": hc.Host,
		}).Info("Checkauth failed at first time retry")
		if flag {
			log.WithFields(log.Fields{
				"raying-cron": hc.Host,
			}).Info("Checkauth success after first time")
		}
		return hc
	}
}

func Monitoring(hc Rayingclient) {
	cronTarget := cron.New()
	spec := "* */5 * * * ?"
	cronTarget.AddFunc(spec, func() {
		go Checkcameraalive(hc)
	})
	cronTarget.Start()
	log.WithFields(log.Fields{
		"raying-cron": "checkalive",
	}).Info("checkalive cron start")
}

func Keepalive(hc Rayingclient) {
	cronTarget := cron.New()
	spec := "*/15 * * * * ?"
	cronTarget.AddFunc(spec, func() {
		go hc.keepalive2(hc.Session)
	})
	cronTarget.Start()
	log.WithFields(log.Fields{
		"raying-cron": "checkalive",
	}).Info("checkalive cron start")
}

func Getgroupid(hc Rayingclient) (bool) {
	var flag bool
	id := 1
	method := "faceGroupManager.getAllGroups"
	obj := map[string]interface{}{
		"id":     id,
		"method": method,
	}
	data, err := hc.Rpc(obj)
	facegroups, err := data.Get("params").Array()
	if err != nil || len(facegroups) == 0 {
		//创建一个人脸库获取groupid
		groupid, err := Creategroupid(hc)
		if err == true {
			hc.GroupID = groupid
			flag = true
		} else {
			hc.GroupID = 0
			flag = false
		}
	}
	if len(facegroups) != 0 {
		GroupID, _ := facegroups[0].(map[string]interface{})["GroupID"].(json.Number).Int()
		hc.GroupID = GroupID
		flag = true
		//return GroupId,true
	}
	write_flag := SaveGroupIDtoml(hc.Host, hc.GroupID)
	if write_flag {
		log.WithFields(log.Fields{
			"raying-cron": hc.Host,
		}).Info("write GroupID success")
	} else {
		log.WithFields(log.Fields{
			"raying-cron": hc.Host,
		}).Info("write GroupID failed")
	}
	return flag
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		Create_Dir(path)
		return PathExists(path)
	}
	return false, err
}
func Create_Dir(path string) (bool) {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Printf("mkdir failed![%v]\n", err)
		return false
	} else {
		fmt.Printf("mkdir success!\n")
		return true
	}
}
func DIR_INIT() () {
	path_list := []string{"./camera_IP", "./camera_GROUPID"}
	for i := 0; i < len(path_list); i++ {
		PathExists(path_list[i])
	}
}
type Facelib struct {
	GroupAlias      string
	Enable          bool
	SearchThreshold float32
	TopRank         int
	GroupType       int
}

func Creategroupid(hc Rayingclient) (int, bool) {
	var GroupId int
	id := 1
	method := "faceGroupManager.createGroup"
	param := Facelib{
		GroupAlias:      "defaultfacelib",
		Enable:          true,
		SearchThreshold: 80,
		TopRank:         1,
		GroupType:       0,
	}
	obj := map[string]interface{}{
		"id":     id,
		"method": method,
		"params": param,
	}
	data, err := hc.Rpc(obj)
	log4go.Info(data)
	facegroups, err := data.Get("params").Get("FaceGroup").Map()
	log4go.Info(facegroups)
	if err == nil {
		id, _ := facegroups["GroupID"].(json.Number).Int()
		GroupId = id
	} else {
		return 0, false
	}
	return GroupId, true
}

func SaveGroupIDtoml(ip string, groupid int) (bool) {
	facegroupuser := Facegroupuser{
		GroupID: groupid,
	}
	//path := "./camera_GROUPID"
	file_path := fmt.Sprintf("./camera_GROUPID/%s_GROUPID.toml", ip)
	//_,err := PathExists(path)
	f, err := os.Create(file_path)
	if err != nil {
		// failed to create/open the file
		log.Fatal(err)
		return false
	}
	if err := toml.NewEncoder(f).Encode(facegroupuser); err != nil {
		// failed to encode
		log.Fatal(err)
		return false
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		log.Fatal(err)
		return false
	}
	return true
}
func ReadGroupIDtoml(ip string) (string, string, int, bool) {
	gp := Facegroupuser{}
	ck := Sessionmaker{}
	path_GI := fmt.Sprintf("./camera_GROUPID/%s_GROUPID.toml", ip)
	path_CK := fmt.Sprintf("./camera_IP/%s_IP.toml", ip)
	if _, err := toml.DecodeFile(path_GI, &gp); err != nil {
		log.Error("read toml file error(%v)", err)
		return "", "", 0, false
	}
	if _, err := toml.DecodeFile(path_CK, &ck); err != nil {
		log.Error("read toml file error(%v)", err)
		return "", "", 0, false
	}
	return ck.Cookies, ck.Authorization, gp.GroupID, true
}


//添加人员信息
func Createphoto(subject model.Person) (*simplejson.Json, error) {

	cookie, authorization, groupid, read_flag := ReadGroupIDtoml(subject.Ip)
	if read_flag != true {
		//读取toml文件失败
		return Datatransform(-100, "GET_Rpc json transform failed", nil), nil
	}

	sub_id := strconv.Itoa(subject.Subject_Id)
	hc := Rayingclient{}
	hc.Init(subject.Ip, cookie, authorization)
	flag := FaceInfoCreateConnect(hc)
	if flag != true {
		log.Error("FaceInfoCreateConnect failed")
		//return Datatransform(-100, "FaceInfoCreateConnect failed", nil), nil
	}
	//flag := true
	id := 1
	method := "faceInfoUpdate.addFace"

	length, filename, err := readfile(subject.Photo)
	if err != nil {
		//图片转换为二进制失败
		return Datatransform(-100, "read photo file failed", nil), nil
	}
	PersonInfo := model.RYPerson_2{
		ID:              sub_id,
		Name:            subject.Name,
		CertificateType: "IC",
		Sex:             "male",
		Country:         "中国",
		City:            "杭州",
		Province:        "浙江",
	}
	ImageInfo := model.RYimg{
		Lengths: [1]int64{length},
		Amount:  1,
	}
	params := model.CreatePhoto{
		GroupID:    groupid,
		PersonInfo: PersonInfo,
		ImageInfo:  ImageInfo,
	}
	body := model.Rpcmodel{
		Id:     id,
		Method: method,
		Params: params,
	}
	files := make(map[string]interface{})
	files["filename"] = filename
	files["photobinary"] = subject.Photo
	files["body"] = body
	ry_subject, err := hc.Rpcp(files)
	if err != nil {
		FaceInfoCreateLooseConnect(hc)
		code, _ := ry_subject.Get("code").Int()
		code_str := strconv.Itoa(code)
		if code == -100 {
			err_msg, _ := ry_subject.Get("err_msg").String()
			return Datatransform(code, err_msg, nil), err
		}
		return Datatransform(code, digital2String(code_str), nil), err
	}
	FaceInfoCreateLooseConnect(hc)
	//PersonInfo_json,_:=Buildjson(PersonInfo)
	//PersonInfo_map,_:=Json2Map(PersonInfo)
	//ry_subject_map,_:=Json2Map(ry_subject)
	koalaperson, _ := koalaadperson(ry_subject, PersonInfo, subject.Photo, subject.Subject_type)
	koalaperson_json, _ := Buildjson(koalaperson)
	return koalaperson_json, err
}

//更新人员信息
func Updatephoto(subject model.Person, face_token string) (*simplejson.Json, error) {

	cookie, authorization, _, read_flag := ReadGroupIDtoml(subject.Ip)
	if read_flag != true {
		//读取toml文件失败
		return Datatransform(-100, "GET_Rpc json transform failed", nil), nil
	}
	hc := Rayingclient{}
	hc.Init(subject.Ip, cookie, authorization)
	flag := FaceInfoCreateConnect(hc)
	if flag != true {
		log.Error("FaceInfoCreateConnect failed")
		//return Datatransform(-100, "FaceInfoCreateConnect failed", nil), nil
	}
	id := 1
	method := "faceInfoUpdate.updateFaceImage"
	_, filename, err := readfile(subject.Photo)
	if err != nil {
		//读取图片信息失败
		return Datatransform(-100, "read photo file failed", nil), nil
	}
	params := map[string]string{
		"FaceToken": face_token,
	}
	body := model.Rpcmodel{
		Id:     id,
		Method: method,
		Params: params,
	}
	files := make(map[string]interface{})
	files["filename"] = filename
	files["photobinary"] = subject.Photo
	files["body"] = body
	ry_subject, err := hc.Rpcp(files)
	if err != nil {
		go FaceInfoCreateLooseConnect(hc)
		code, _ := ry_subject.Get("code").Int()
		code_str := strconv.Itoa(code)
		if code == -100 {
			err_msg, _ := ry_subject.Get("err_msg").String()
			return Datatransform(code, err_msg, nil), err
		}
		return Datatransform(code, digital2String(code_str), nil), err
	}
	go FaceInfoCreateLooseConnect(hc)

	return Getfacetoken(ry_subject), err
}

func Updateperson(subject model.Person, person_id int, facetoken string) (*simplejson.Json, error) {

	sub_id := strconv.Itoa(subject.Subject_Id)
	cookie, authorization, _, read_flag := ReadGroupIDtoml(subject.Ip)
	if read_flag != true {
		//读取toml文件失败
		return Datatransform(-100, "read toml file failed", nil), nil
	}

	hc := Rayingclient{}
	hc.Init(subject.Ip, cookie, authorization)
	flag := FaceInfoCreateConnect(hc)
	if flag != true {
		log.Error("FaceInfoCreateConnect failed")
		//return Datatransform(-100, "FaceInfoCreateConnect failed", nil), nil
	}
	id := 1
	method := "faceInfoUpdate.updatePersonInfo"
	PersonInfo_obj := model.RYPerson_2{
		ID:              sub_id,
		Name:            subject.Name,
		CertificateType: "IC",
		Sex:             "male",
		Country:         "中国",
		City:            "杭州",
		Province:        "浙江",
	}
	fmt.Println(PersonInfo_obj)
	PersonInfo := map[string]interface{}{
		"Name":            subject.Name,
		"CertificateType": "IC",
		"Sex":             "male",
		"Country":         "中国",
		"City":            "杭州",
		"Province":        "浙江",
	}
	params := map[string]interface{}{
		"PersonID":   person_id,
		"PersonInfo": PersonInfo,
	}
	files := map[string]interface{}{
		"id":     id,
		"method": method,
		"params": params,
	}

	ry_subject, err := hc.Rpc(files)
	if err != nil {
		go FaceInfoCreateLooseConnect(hc)
		code, _ := ry_subject.Get("code").Int()
		code_str := strconv.Itoa(code)
		if code == -100 {
			err_msg, _ := ry_subject.Get("err_msg").String()
			return Datatransform(code, err_msg, nil), err
		}
		return Datatransform(code, digital2String(code_str), nil), err
	}
	go FaceInfoCreateLooseConnect(hc)

	facetoken_list := [1]string{facetoken}
	FaceToken := map[string]interface{}{
		"FaceToken": facetoken_list,
	}
	mokemap := map[string]interface{}{
		"id":     1,
		"params": FaceToken,
	}
	mokemap_json, _ := Buildjson(mokemap)
	koalaperson, _ := koalaadperson(mokemap_json, PersonInfo_obj, subject.Photo, subject.Subject_type)
	koalaperson_json, _ := Buildjson(koalaperson)
	return koalaperson_json, nil

}

func Deleteperson(subject model.Person) (*simplejson.Json, error) {

	sub_id := strconv.Itoa(subject.Subject_Id)
	cookie, authorization, _, read_flag := ReadGroupIDtoml(subject.Ip)
	if read_flag != true {
		//读取toml文件失败
		return Datatransform(-100, "read toml file failed", nil), nil
	}

	hc := Rayingclient{}
	hc.Init(subject.Ip, cookie, authorization)
	flag := FaceInfoCreateConnect(hc)
	if flag != true {
		log.Error("FaceInfoCreateConnect failed")
		//return Datatransform(-100, "FaceInfoCreateConnect failed", nil), nil
	}
	id := 1
	method := "faceInfoUpdate.deletePersonByGuid"
	body := map[string]interface{}{
		"CertificateType": "IC",
		"ID":              sub_id,
	}
	data := map[string]interface{}{
		"id":     id,
		"method": method,
		"params": body,
	}
	ry_subject, err := hc.Rpc(data)
	if err != nil {
		go FaceInfoCreateLooseConnect(hc)
		code, _ := ry_subject.Get("code").Int()
		code_str := strconv.Itoa(code)
		if code == -100 {
			err_msg, _ := ry_subject.Get("err_msg").String()
			return Datatransform(code, err_msg, nil), err
		}
		return Datatransform(code, digital2String(code_str), nil), err
	}
	go FaceInfoCreateLooseConnect(hc)
	return ry_subject, err

}

func GetPersonInfoFromID(subject model.Person, ID string) (*simplejson.Json, error) {

	cookie, authorization, _, read_flag := ReadGroupIDtoml(subject.Ip)
	if read_flag != true {
		//读取toml文件失败
		return Datatransform(-100, "read toml file failed", nil), nil
	}
	hc := Rayingclient{}
	hc.Init(subject.Ip, cookie, authorization)
	flag := FaceInfoFindCreateConnect(hc)
	if flag == 0 {
		//获取实例失败
		return Datatransform(-100, "获取实例失败", nil), nil
	}
	id := 1
	method := "faceInfoFind.getPersonInfoByID"
	params := map[string]string{
		"CertificateType": "IC",
		"ID":              ID,
	}
	body := model.Rpcmodel_findperson{
		Id:     id,
		Method: method,
		Params: params,
		Object: flag,
	}
	ry_subject, err := hc.Rpc(body)
	if err != nil {
		go FaceInfoFindLooseConnect(hc, flag)
		code, _ := ry_subject.Get("code").Int()
		code_str := strconv.Itoa(code)
		if code == -100 {
			err_msg, _ := ry_subject.Get("err_msg").String()
			return Datatransform(code, err_msg, nil), err
		}
		return Datatransform(code, digital2String(code_str), nil), err
	}
	go FaceInfoFindLooseConnect(hc, flag)
	return ry_subject.Get("params"), nil

}

func FaceInfoCreateConnect(hc Rayingclient) (bool) {

	id := 1
	method := "faceInfoUpdate.create"
	obj := map[string]interface{}{
		"id":     id,
		"method": method,
	}
	data, err := hc.Rpc(obj)
	log4go.Info(data)
	if data == nil {
		return false
	}
	_, ok := data.CheckGet("result")
	if ok != true {
		return false
	}

	result := data.Get("result")
	flag, err := result.Bool()
	if err != nil {
		return false
	}
	if flag {
		return flag
	} else {
		return flag
	}
}

func FaceInfoCreateLooseConnect(hc Rayingclient) (bool) {

	id := 1
	method := "faceInfoUpdate.destroy"

	obj := map[string]interface{}{
		"id":     id,
		"method": method,
	}
	data, err := hc.Rpc(obj)
	_, ok := data.CheckGet("result")
	if ok != true {
		return false
	}
	result := data.Get("result")
	flag, err := result.Bool()
	if err != nil {
		return false
	}
	if flag {
		return flag
	} else {
		return flag
	}
}

func file2Bytes(filename string) ([]byte, int, string, error) {
	// File
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, "", err
	}
	defer file.Close()
	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return nil, 0, "", err
	}
	// []byte
	data := make([]byte, stats.Size())
	count, err := file.Read(data)
	if err != nil {
		return nil, 0, "", err
	}
	fmt.Printf("read file %s len: %d \n", filename, count)
	return data, count, file.Name(), nil
}

func readfile(filename string) (int64, string, error) {
	// File
	file_size, err := os.Stat(filename)
	//file, err := os.Open(filename)

	if err != nil {
		return 0, "", err
	}
	//defer file.Close()
	// FileInfo:
	//stats, err := file.Stat()
	//if err != nil {
	//	return nil,0,"", err
	//}
	// []byte
	//data := make([]byte, stats.Size())
	//count, err := file.Read(data)
	//if err != nil {
	//	return nil,0,"", err
	//}
	fmt.Printf("read file %s size: %d \n", filename, file_size.Size())
	return file_size.Size(), file_size.Name(), nil
}

func Createphoto2(subject model.Person) () {
	cookie, authorization, _, read_flag := ReadGroupIDtoml(subject.Ip)
	if read_flag != true {
		//读取toml文件失败
	}
	hc := Rayingclient{}
	hc.Init(subject.Ip, cookie, authorization)

	go hc.keepalive(cookie)
}

func FaceInfoFindCreateConnect(hc Rayingclient) (int) {

	id := 1
	method := "faceInfoFind.create"

	obj := map[string]interface{}{
		"id":     id,
		"method": method,
	}
	data, err := hc.Rpc(obj)
	if data == nil {
		return 0
	}
	result := data.Get("result")
	flag, err := result.Int()
	if err != nil {
		return 0
	}
	if flag != 0 {
		return flag
	} else {
		return 0
	}
}

func FaceInfoFindLooseConnect(hc Rayingclient, object int) (bool) {

	id := 1
	method := "faceInfoFind.destroy"
	obj := map[string]interface{}{
		"id":     id,
		"method": method,
		"object": object,
	}
	data, err := hc.Rpc(obj)
	result := data.Get("result")
	flag, err := result.Bool()
	if err != nil {
		return false
	}
	if flag {
		return flag
	} else {
		return flag
	}
}

func Getevent(Ip, callback string) (*simplejson.Json, error) {

	//var(
	//	ry_subject *simplejson.Json
	//	err error
	//)

	cookie, authorization, _, read_flag := ReadGroupIDtoml(Ip)
	if read_flag != true {
		//读取toml文件失败
		return Datatransform(-100, "read toml file failed", nil), nil
	}
	fmt.Println(callback)
	Callback_url = callback
	hc := Rayingclient{}
	hc.Init(Ip, cookie, authorization)

	ry_subject, err := hc.Rpcp_event()

	//hc.startWebsocket()

	if err != nil {
		//go FaceInfoCreateLooseConnect(hc)
		code, _ := ry_subject.Get("code").Int()
		code_str := strconv.Itoa(code)
		if code == -100 {
			err_msg, _ := ry_subject.Get("err_msg").String()
			return Datatransform(code, err_msg, nil), err
		}
		return Datatransform(code, digital2String(code_str), nil), err
	}
	//go FaceInfoCreateLooseConnect(hc)

	//facetoken_list := [1]string{facetoken}
	//FaceToken := map[string]interface{}{
	//	"FaceToken": facetoken_list,
	//}
	//mokemap := map[string]interface{}{
	//	"id":     1,
	//	"params": FaceToken,
	//}
	//mokemap_json, _ := Buildjson(mokemap)
	//koalaperson, _ := koalaadperson(mokemap_json, PersonInfo_obj, subject.Photo, subject.Subject_type)
	//koalaperson_json, _ := Buildjson(koalaperson)
	return ry_subject, nil

}

func WriteWithBufio(name, content string) {
	if fileObj, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err == nil {
		defer fileObj.Close()
		writeObj := bufio.NewWriterSize(fileObj, 4096)
		//
		if _, err := writeObj.WriteString(content); err == nil {
			fmt.Println("Successful appending buffer and flush to file with bufio's Writer obj WriteString method", content)
		}

		//使用Write方法,需要使用Writer对象的Flush方法将buffer中的数据刷到磁盘
		buf := []byte(content)
		if _, err := writeObj.Write(buf); err == nil {
			fmt.Println("Successful appending to the buffer with os.OpenFile and bufio's Writer obj Write method.", content)
			if err := writeObj.Flush(); err != nil {
				panic(err)
			}
			fmt.Println("Successful flush the buffer data to file ", content)
		}
	}
}
