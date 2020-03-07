package service

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/snowflake"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"zhiyuan/device_server/raying_api/internal/model"
)

type CreateId struct {
}

func (c *CreateId) CreateID() string {
	sf, err := snowflake.New(1)
	if err != nil {
		panic(err)
	}
	uuid, _ := sf.Generate()
	str_uuid := strconv.FormatUint(uuid, 10)
	fmt.Println(str_uuid)
	return str_uuid
}

type Rayingclient struct {
	Host string
	//Ip	string
	Username      string
	Password      string
	Session       string
	GroupID       int
	IsConnected   bool
	Authorization string
	//PostParams model.Rpcmodel
	Result model.NormalResponsemodel
}

func (s *Rayingclient) Init(ip, session, authorization string) {
	s.Host = ip
	s.Username = Username
	s.Password = Password
	s.Session = session
	s.Authorization = authorization
	//s.Ip = Base_ip
}
func DoResponse(resp *http.Response) (*simplejson.Json, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	//io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log4go.Debug(err)
	}
	log4go.Debug(body)
	log4go.Debug(string(body))
	//io.Copy(ioutil.Discard,resp.Body)
	if err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("Read response body error")
	}
	log4go.Debug(string(body))

	jdata, err := simplejson.NewJson(body)
	if err != nil {
		log4go.Error(err.Error())
		if strings.Index(string(body), "<span>记住我</span>") == -1 {
			return nil, errors.New("Face++返回报文错误")
		} else {
			return nil, errors.New("登陆失效，请重新登录!")
		}

	}

	code, _ := jdata.Get("code").Int()
	if code != 0 {
		desc, _ := jdata.Get("result").String()
		log4go.Error(desc)
		return nil, errors.New(desc)
	}
	return jdata, nil
}

func (s *Rayingclient) CheckauthFirsttime(host string) (string, string) {
	baseUrl := Http + host
	HttpClient := &http.Client{
	}
	id := 1
	method := "faceGroupManager.getAllGroups"
	param := map[string]interface{}{
		"id":     id,
		"method": method,
	}
	jsonValue, err := json.Marshal(param)
	req, err := http.NewRequest(MethodPost, baseUrl+Raying_api, bytes.NewBuffer(jsonValue))

	if err != nil {
		log4go.Error(err)
	}
	resp, err := HttpClient.Do(req)
	Authenticate := resp.Header.Get("WWW-Authenticate")
	setcookie := resp.Header.Get("Set-Cookie")
	if Authenticate == "" || setcookie == "" {
		return Authenticate, setcookie
	}
	params := strings.Split(Authenticate, ",")
	redult := s.transmap(params)
	realm := redult["Digest realm"]
	qop := redult["qop"]
	nonce := redult["nonce"]
	opaque := redult["opaque"]
	cnonce := "ksjdfljwofsldj4687skjd"
	nc := "00000001"
	A1 := s.Username + ":" + realm + ":" + s.Password
	//A2 := MethodGet+":"+Digest_url
	A2 := MethodPost + ":" + Raying_api
	HA1 := s.Cal_md5(A1)
	HA2 := s.Cal_md5(A2)
	response := s.Cal_md5(HA1 + ":" + nonce + ":" + nc + ":" + cnonce + ":" + qop + ":" + HA2)
	authorization_str := fmt.Sprintf("Digest username=\"%s\",realm=\"%s\",nonce=\"%s\",uri=\"%s\",qop=auth,cnonce=\"%s\",response=\"%s\",opaque=\"%s\",nc=\"%s\"", s.Username, realm, nonce, Raying_api, cnonce, response, opaque, nc)
	//authorization_str := fmt.Sprintf("Digest username=\"%s\",realm=\"%s\",nonce=\"%s\",uri=\"%s\",qop=\"%s\",cnonce=\"%s\",response=\"%s\",opaque=\"%s\",nc=\"%s\"", s.Username, realm,nonce, Digest_url, qop,cnonce, response, opaque,nc)
	authorization := authorization_str
	cookie := strings.Split(setcookie, "=")
	return authorization, cookie[1]
}

func (s *Rayingclient) Checkauth() (bool) {
	s.IsConnected = false
	host := s.Host
	authorization, cookies := s.CheckauthFirsttime(host)
	baseUrl := Http + host
	id := 1
	method := "faceGroupManager.getAllGroups"
	param := map[string]interface{}{
		"id":     id,
		"method": method,
	}
	jsonValue, err := json.Marshal(param)
	req, err := http.NewRequest(MethodPost, baseUrl+Raying_api, bytes.NewBuffer(jsonValue))
	//HttpClient := &http.Client{Timeout:30*time.Second}
	HttpClient := &http.Client{}
	//req, err := http.NewRequest(MethodGet, baseUrl+Digest_url, nil)
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", cookies))
	resp, err := HttpClient.Do(req)

	if err != nil {
		log4go.Error(err)
		return false
	}
	if resp == nil {
		return false
	}
	if resp.StatusCode == 200 {
		log.WithFields(log.Fields{
			"raying": s.Host,
		}).Info("connection successed do keepalive next")
		log.WithFields(log.Fields{
			"raying": "resp",
		}).Info(resp)
		go s.keepalive(cookies)
		go s.Rpcp_event()

	} else {
		log.WithFields(log.Fields{
			"raying": s.Host,
		}).Info("connection failed please retry")
		s.IsConnected = false
		return s.IsConnected
	}
	s.Session = cookies
	s.Authorization = authorization
	s.IsConnected = true
	flag := s.Savetoml(s.Host, s.Session, s.Authorization)
	if flag {
		log.WithFields(log.Fields{
			"raying": s.Host,
		}).Info("save2toml success")
	} else {
		log.WithFields(log.Fields{
			"raying": s.Host,
		}).Info("save2toml failed")
	}
	return s.IsConnected
}
func (s *Rayingclient) keepalive(session string) {
	baseUrl := Http + s.Host
	data := url.Values{}
	HttpClient := &http.Client{}
	data.Set("username", "")
	keepalive_url := fmt.Sprintf("/api/common/keepalive?session=%s&active=true", session)
	for {
		time.Sleep(15 * time.Second)
		req, err := http.NewRequest("POST", baseUrl+keepalive_url, bytes.NewBufferString(data.Encode()))
		if err != nil {
			log4go.Error(err)
		}
		resp, err := HttpClient.Do(req)
		if resp == nil {
			log4go.Info(resp)
		} else {
			if resp.StatusCode == 200 {
				//log4go.Info(resp.StatusCode)
				continue
			} else {
				log4go.Error(resp.StatusCode)
				continue
			}
		}
	}
}
func (s *Rayingclient) keepalive2(session string) {
	baseUrl := Http + s.Host
	data := url.Values{}
	HttpClient := &http.Client{}
	data.Set("username", "")
	keepalive_url := fmt.Sprintf("/api/common/keepalive?session=%s&active=true", session)
	req, err := http.NewRequest("POST", baseUrl+keepalive_url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log4go.Error(err)
	}
	resp, err := HttpClient.Do(req)
	if resp.StatusCode == 200 {
		log4go.Info(resp.StatusCode)
	} else {
		log4go.Error(resp.StatusCode)
	}
	//return
}

func (s *Rayingclient) Savetoml(ip, cookie, authorization string) (bool) {
	sessionsave := Sessionmaker{
		Cookies:       cookie,
		Authorization: authorization,
	}
	f, err := os.Create(fmt.Sprintf("./camera_IP/%s_IP.toml", ip))
	if err != nil {
		// failed to create/open the file
		log.Fatal(err)
		return false
	}
	if err := toml.NewEncoder(f).Encode(sessionsave); err != nil {
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

func (s *Rayingclient) Cal_md5(str string) (string) {
	data := []byte(str)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
func (s *Rayingclient) transmap(str []string) (map[string]string) {
	res := make(map[string]string)
	for _, v := range str {
		params := strings.Split(v, "=")
		if len(params) == 2 {
			res[strings.TrimSpace(params[0])] = strings.TrimSpace(params[1][1 : len(params[1])-1])
		}
	}
	return res
}

func (s *Rayingclient) Rpc(param interface{}) (*simplejson.Json, error) {
	//POST
	baseUrl := Http + s.Host
	jsonValue, err := json.Marshal(param)
	if err != nil {
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc json transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}

	req, err := http.NewRequest(MethodPost, baseUrl+Raying_api, bytes.NewBuffer(jsonValue))
	req.Close = true
	req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", s.Session))
	req.Header.Add("Authorization", s.Authorization)

	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc http NewRequest transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	HttpClient := &http.Client{}
	resp, err := HttpClient.Do(req)
	if resp == nil{
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc is nil"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	if resp.StatusCode == 401{
		s.Checkauth()
		req, _ := http.NewRequest(MethodPost, baseUrl+Raying_api, bytes.NewBuffer(jsonValue))
		req.Close = true
		req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", s.Session))
		req.Header.Add("Authorization", s.Authorization)
	}

	//fmt.Println(resp)
	if err != nil {
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc http client transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	data, err := DoResponse(resp)
	if data == nil {
		s.Result.Code = -100
		s.Result.Message = "连接中断请稍后重连"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, errors.New("RPC FAILED")
	}
	ok, err := data.Get("result").Bool()
	if err != nil {
		_, err := data.Get("result").Int()
		if err == nil {
			return data, nil
		}
	}
	if ok {
		return data, nil
	} else {
		s.Result.Code, _ = data.Get("error").Get("code").Int()
		s.Result.Message, _ = data.Get("error").Get("message").String()
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, errors.New("RPC FAILED")
	}
}
func (s *Rayingclient) Rpcp(param map[string]interface{}) (*simplejson.Json, error) {
	//POST
	baseUrl := Http + s.Host
	filename := param["filename"]
	photobinary := param["photobinary"]
	bodys := param["body"]
	file, err := os.Open(photobinary.(string))
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	boundary := writer.Boundary()
	writer.SetBoundary(boundary)
	jsonBytes, err := json.Marshal(bodys)
	_ = writer.WriteField("json", string(jsonBytes), "application/json")
	part, err := writer.CreateFormFile("data", filename.(string), "image/jpeg")

	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "RPCP 组成JSON失败"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "RPCP 表单文件拷贝失败"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}

	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc json transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	HttpClient := &http.Client{}
	req, err := http.NewRequest(MethodPost, baseUrl+Raying_api, body)
	req.Close = true
	req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", s.Session))
	req.Header.Set("connection", "close")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc http NewRequest transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	resp, err := HttpClient.Do(req)

	if resp == nil{
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpcp is nil"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}

	if resp.StatusCode == 401{
		s.Checkauth()
		req, _ := http.NewRequest(MethodPost, baseUrl+Raying_api, body)
		req.Close = true
		req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", s.Session))
		req.Header.Set("connection", "close")
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc http client transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	data, err := DoResponse(resp)
	if data == nil {
		s.Result.Code = -100
		s.Result.Message = "连接中断请稍后重连"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, errors.New("RPCP FAILED")
	}
	ok, _ := data.Get("result").Bool()
	if ok {
		return data, nil
	} else {
		s.Result.Code, _ = data.Get("error").Get("code").Int()
		s.Result.Message, _ = data.Get("error").Get("message").String()
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, errors.New("RPCP FAILED")
	}
}

func cal_md5(str string) (string) {
	data := []byte(str)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func CreateConnection(ip, session, authorization string) (Rayingclient) {
	hc := Rayingclient{}
	hc.Init(ip, session, authorization)
	return hc
}

func Checkcameraalive(hc Rayingclient) {
	out, _ := exec.Command("ping", hc.Host, "-c 5", "-i 3", "-w 10").Output()
	if strings.Contains(string(out), "Destination Host Unreachable") {
		//fmt.Println("TANGO DOWN")
		log.WithFields(log.Fields{
			"raying": "lose connected",
		}).Info("TARGET SHUIT DOWN TRY RECONNECTIVE")
		go Retry(10000, 30, hc.Checkauth)
		//return false
	} else {
		log.WithFields(log.Fields{
			"raying": "connected",
		}).Info("IT'S ALIVEEE")
		//return true
	}
}

func (s *Rayingclient) Rpcp_event() (*simplejson.Json, error) {
	//POST
	callback := "127.0.0.1:8081/event/event"
	Callback_url =	callback
	baseUrl := Http + s.Host
	req, err := http.NewRequest(MethodPost, baseUrl+Raying_event, nil)
	req.Header.Add("Cookie", fmt.Sprintf("SessionID=%s", s.Session))
	if err != nil {
		log4go.Error(err)
		s.Result.Code = -100
		s.Result.Message = "GET_Rpc http NewRequest transform failed"
		s.Result.Data = nil
		jsondata, _ := Buildjson(s.Result)
		return jsondata, err
	}
	HttpClient := &http.Client{}
	resp, err := HttpClient.Do(req)

	eof := []byte{}
	EOF := append(eof, 255, 217)
	buf := make([]byte, 150000) // any non zero value will do, try '1'.
	data_list := make([]RecordsPKG, 0)
	tempchan := make(chan RecordsPKG, 1)
	for {
		//buf := make([]byte, 150000)
		time.Sleep(1 * time.Second)
		n, err := resp.Body.Read(buf)
		if n == 0 && err != nil { // simplified
			continue
		}
		body := buf[:n]
		if len(body) == 0{
			continue
		}
		if len(data_list) == 0 {
			if bytes.Compare(body[n-2:n], EOF) == 0 {
				log4go.Info("send at first time")
				RG_OBJ := RecordsPKG{}
				RG_OBJ.Dismantling(buf[:n])
				go RG_OBJ.Callback(s.Host)
			} else {
				log4go.Info("pkg length more than 150000")
				RG_OBJ := RecordsPKG{}
				RG_OBJ.Dismantling(buf[:n])
				log4go.Info(RG_OBJ.Json_data)
				log4go.Info(len(RG_OBJ.Json_data))
				//log4go.Info(RG_OBJ.Json_data_byte)
				//log4go.Info(RG_OBJ.Photo_data)
				//oo:=Carry(RG_OBJ)
				//log4go.Info(oo)
				tempchan <- RG_OBJ
				data_list = append(data_list, RG_OBJ)
			}
		} else {
			if bytes.Compare(body[n-2:n], EOF) == 0 {
				log4go.Info("send at next time")
				RG_OBJ := RecordsPKG{}
				RG_OBJ.Dismantling(buf[:n])
				//temp := data_list[0]
				task := <-tempchan
				log4go.Info(task.Json_data)
				log4go.Info(task.Json_data_byte)
				log4go.Info(len(task.Json_data))
				Receive(task, RG_OBJ, s.Host)
				data_list = make([]RecordsPKG, 0)
			}
		}
	}
	return nil, nil
}

type RecordsPKG struct {
	Photo_data         []byte
	Jsoncontent_length int
	Json_data          string
	Json_data_byte		[]byte
	Pkg_length         int
	Callback_url       string
	Ip                 string
}

func Carry(target RecordsPKG) func() RecordsPKG {
	obj := RecordsPKG{}
	return func() RecordsPKG {
		obj = target
		return obj
	}
}

func (RG *RecordsPKG) Dismantling(body []byte) () {

	var (
		photo_data         []byte
		jsoncontent_length []byte
	)

	json_str_first := "{\"method\":\"pictureEvent.notify\",\"params\":"
	Content_Length := "Content-Length: "
	Sequence := "Sequence: "
	END := "}}}"
	json_str_first_byte := []byte(json_str_first)
	content_length := []byte(Content_Length)
	end := []byte(END)
	sequence := []byte(Sequence)

	index_json_str_first_byte := bytes.Index(body, json_str_first_byte)
	//index_content_length := bytes.Index(body,content_length)
	indexd_end := bytes.Index(body, end)
	indexe_sequence := bytes.Index(body, sequence)
	indexe_sequence_last := bytes.LastIndex(body, sequence)
	index_content_length_last := bytes.LastIndex(body, content_length)
	//JSON_DATA
	if index_json_str_first_byte != -1 && indexd_end != -1 {
		jsonstr := body[index_json_str_first_byte : indexd_end+3]
		RG.Json_data_byte = jsonstr
		RG.Json_data = String(jsonstr)
		log4go.Info(RG.Json_data)
	}
	//content_length
	if indexe_sequence_last != -1 && index_content_length_last != -1 {
		jsoncontent_length = body[index_content_length_last+16 : indexe_sequence_last-2]
		str := String(jsoncontent_length)
		content_length, _ := strconv.Atoi(str)
		RG.Jsoncontent_length = content_length
		//log4go.Info(str)
	}
	//PHOTO DATA
	if indexe_sequence != -1 && index_json_str_first_byte != -1 {
		sequence_number := body[indexe_sequence+10 : index_json_str_first_byte]
		photo_data = body[indexe_sequence_last+len(sequence_number)+len(Sequence):]
		RG.Photo_data = photo_data
		RG.Pkg_length = len(photo_data)
		//str := String(photo_data)
		//log4go.Info(len(photo_data))
		//log4go.Info(str)
	} else {
		RG.Photo_data = body
		RG.Json_data = ""
		RG.Jsoncontent_length = len(body)
		RG.Pkg_length = len(body)
	}
	//log4go.Info(RG.Jsoncontent_length)
	//log4go.Info(RG.Pkg_length)
	return
}
func Receive(obj1, obj2 RecordsPKG, ip string) () {
	log4go.Info("begin to merge pkg")
	if obj1.Pkg_length+obj2.Pkg_length == obj1.Jsoncontent_length {
		//log4go.Info("------------------------------------receive")
		//log4go.Info(obj1.Json_data_byte)
		//log4go.Info("------------------------------------photobyte")
		//log4go.Info(obj1.Photo_data)
		for _, v := range obj2.Photo_data {
			obj1.Photo_data = append(obj1.Photo_data, v)
		}
		log4go.Info("merge pkg success")
		obj1.Json_data = String(obj1.Json_data_byte)
		log4go.Info(obj1.Json_data)
		log4go.Info(len(obj1.Photo_data))
		//go obj1.Callback(ip)
	} else {
		log4go.Info("merge pkg failed")
		//return RecordsPKG{Pkg_length:0,Jsoncontent_length:0,Json_data:"",Photo_data:nil,},false
	}
	//log4go.Info(len(obj1.Photo_data))
	//log4go.Info(obj1.Jsoncontent_length)
	//log4go.Info(obj1.Pkg_length)
	//return obj1,true
}

func (RG *RecordsPKG) Callback(ip string) () {

	log4go.Info(RG.Json_data)
	event_obj := model.EventPerson{}
	//json str 转map
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(RG.Json_data), &dat); err != nil {
		log4go.Error("json转map失败")
	}
	RYEVENT_JSON, _ := Buildjson(dat)
	Events, _ := RYEVENT_JSON.Get("params").Get("EventInfo").Get("Events").Array()
	for _, v := range Events {
		//log4go.Info(v.(map[string]interface{})["RecognizeResults"])
		RecognizeResults := v.(map[string]interface{})["RecognizeResults"].([]interface{})
		for _, va := range RecognizeResults {
			log4go.Info(va)
			event_obj.PersonName = va.(map[string]interface{})["PersonInfo"].(map[string]interface{})["Name"].(string)
			event_obj.SubjectId = va.(map[string]interface{})["PersonInfo"].(map[string]interface{})["ID"].(string)
			event_obj.PersonType = "1"
			event_obj.Similarity = va.(map[string]interface{})["SearchScore"].(json.Number).String()
		}
		//log4go.Info(v.(map[string]interface{})["UTC"].(int64))
		UTC_time, _ := v.(map[string]interface{})["UTC"].(json.Number).Int64()
		event_obj.Timestamp = UTC_time - (8 * 60 * 60)
	}
	encodeString := base64.StdEncoding.EncodeToString(RG.Photo_data)
	event_obj.FacePicture = encodeString
	event_obj.CameraIp = ip
	log4go.Info("send msg send msg")
	json_obj := map[string]interface{}{
		"personName":  event_obj.PersonName,
		"subjectId":   event_obj.SubjectId,
		"personType":  event_obj.PersonType,
		"facePicture": event_obj.FacePicture,
		"timestamp":   event_obj.Timestamp,
		"similarity":  event_obj.Similarity,
		"cameraIp":    event_obj.CameraIp,
	}
	RG.Rpc(json_obj)
	//log4go.Info("send msg send msg")
}

func (RG *RecordsPKG) Rpc(param interface{}) () {
	//POST
	RG.Callback_url = Callback_url
	baseUrl := Http + RG.Callback_url
	log4go.Info(baseUrl)
	jsonValue, err := json.Marshal(param)
	if err != nil {
	}
	log4go.Info(string(jsonValue))
	req, err := http.NewRequest(MethodPost, baseUrl, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Close = true
	log4go.Info(req)
	if err != nil {
		log4go.Error(err)
	}
	HttpClient := &http.Client{}
	resp, err := HttpClient.Do(req)
	if err != nil {

	}
	//data, err := DoResponse(resp)
	if resp ==nil{
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	data := make(map[string]interface{})
	if err != nil {
		log4go.Error(err.Error())
		//return nil, err
	}
	err1 := json.Unmarshal(body, &data)
	if err1 != nil {
		log4go.Error(err1.Error())
		//return nil, err
	}
	//if data != nil {
	log4go.Info("send msg success")
	log4go.Info(data)
	//}

}

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}
func Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
