// koala api
// 访问考拉服务器
package koala

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"net/http/cookiejar"

	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
)

var baseUrl = ""
var koalaHost = ""
var jar, _ = cookiejar.New(nil)


//func InitLogin(){
//
//}


func Init(url string) {
	//baseUrl = "http://" + config.Gconf.KoalaHost + ":" + strconv.Itoa(config.Gconf.KoalaPort)
	//baseUrl = "http://"+url
	baseUrl = "http://hz91zo.oicp.vip:10880"
	koalaHost = "hz91zo.oicp.vip:10880"
	//koalaHost = config.Gconf.KoalaHost + ":" + strconv.Itoa(config.Gconf.KoalaPort)
	log4go.Info("koala url: " + baseUrl)
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	//jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	/*
		var err error
		jar, err = cookiejar.New(nil)
		if err != nil {
			log4go.Crash(err)
		}

		if err := KoalaLogin(config.Gconf.KoalaUsername, config.Gconf.KoalaPassword); err != nil {
			log4go.Crash(err)
		}
	*/

}

//func doResponse(body *[]byte) error {
func doResponse(resp *http.Response) (*simplejson.Json, error) {
	log4go.Debug(resp.Status)
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("Read response body error")
	}
	log4go.Debug(string(body))

	jdata, err := simplejson.NewJson(body)
	if err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("Face++返回报文错误")
	}

	code, _ := jdata.Get("code").Int()
	if code != 0 {
		desc, _ := jdata.Get("desc").String()
		log4go.Error(desc)
		return nil, errors.New(desc)
	}
	return jdata, nil
}

func KoalaLogin(username string, password string) error {
	client := &http.Client{
		Jar: jar,
	}

	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)

	req, err := http.NewRequest("POST", baseUrl+"/auth/login", bytes.NewBufferString(data.Encode()))
	if err != nil {
		log4go.Error(err)
		return errors.New("New request error")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Koala Admin")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return err
	}

	_, err = doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return err
	}

	return nil

}

func AddPhoto(photo *multipart.File) (int, error) {
	client := &http.Client{
		Jar: jar,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", "photo.jpg","image/jpeg")
	if err != nil {
		return -1, err
	}
	_, err = io.Copy(part, *photo)
	if err != nil {
		log4go.Error(err.Error())
		return -1, err
	}

	//writer.WriteField("aa", "aa")

	err = writer.Close()
	if err != nil {
		return -1, err
	}

	request, err := http.NewRequest("POST", baseUrl+"/subject/photo", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	log4go.Debug(request.Header)

	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return -1, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return -1, err
	}

	photo_id, _ := resp_json.Get("data").Get("id").Int()

	return photo_id, nil
}

func AddSubject(params *map[string]interface{}, photo *multipart.File) (*map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 上传底库照片
	photo_id, err := AddPhoto(photo)
	if err != nil {
		log4go.Error(err)
		return nil, err
	}

	jsdata := simplejson.New()
	for key, val := range *params {
		jsdata.Set(key, val)
	}
	photo_ids := []int{photo_id}
	jsdata.Set("photo_ids", photo_ids)
	byte_data, _ := jsdata.MarshalJSON()
	log4go.Debug(string(byte_data))

	// 新增人员
	req, err := http.NewRequest("POST", baseUrl+"/subject", strings.NewReader(string(byte_data)))
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	var res_data = make(map[string]interface{})
	res_data["id"] = resp_json.Get("data").Get("id")
	res_data["photo_id"] = photo_id
	res_data["name"] = resp_json.Get("data").Get("name")

	return &res_data, nil
}

// 修改subject
func ModSubject(params *map[string]interface{}, photo *multipart.File) (*map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 上传底库照片
	photo_id, err := AddPhoto(photo)
	if err != nil {
		log4go.Error(err)
		return nil, err
	}

	jsdata := simplejson.New()
	for key, val := range *params {
		jsdata.Set(key, val)
	}
	photo_ids := []int{photo_id}
	jsdata.Set("photo_ids", photo_ids)
	byte_data, _ := jsdata.MarshalJSON()
	log4go.Debug(string(byte_data))

	// 修改人员
	subject_id := (*params)["subject_id"].(int)
	req, err := http.NewRequest("PUT", baseUrl+"/subject/"+strconv.Itoa(subject_id), strings.NewReader(string(byte_data)))
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	var res_data = make(map[string]interface{})
	res_data["id"] = resp_json.Get("data").Get("id")
	res_data["photo_id"] = photo_id
	res_data["name"] = resp_json.Get("data").Get("name")

	return &res_data, nil
}

func DeleteSubject(subject_id int) error {
	client := &http.Client{
		Jar: jar,
	}

	// 新增人员
	req, err := http.NewRequest("DELETE", baseUrl+"/subject/"+strconv.Itoa(subject_id), nil)
	if err != nil {
		log4go.Error(err)
		return errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return err
	}

	_, err = doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return err
	}

	return nil
}

func GetSubjects(category string) (*simplejson.Json, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 新增人员
	req, err := http.NewRequest("GET", baseUrl+"/mobile-admin/subjects/list?category="+category+"&size=1000", nil)
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	return resp_json, nil
}

func GetDisplayDevice(device_token string) (*simplejson.Json, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 新增人员
	req, err := http.NewRequest("GET", baseUrl+"/screen/get-display-config?device_token="+device_token, nil)
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}

	return resp_json, nil
}
