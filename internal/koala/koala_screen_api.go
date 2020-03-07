package koala


// koala screen api
// 访问考拉服务器 -- 迎宾部分的api

import (
//"bytes"
//"config"
//"errors"
//"io"
//"io/ioutil"
//"mime/multipart"
//"net/http"
"net/http/httputil"
"net/url"
//"strconv"
//"strings"

//"net/http/cookiejar"

"github.com/alecthomas/log4go"
//"github.com/bitly/go-simplejson"
//"gopkg.in/gin-gonic/gin.v1"
"github.com/gin-gonic/gin"
)

func ReverseProxy2(c *gin.Context) error {

	//target := baseUrl + c.Request.RequestURI

	/*
		director := func(req *http.Request) {
			r := c.Request
			req = r
			log4go.Debug(req.URL.Scheme)
			req.URL.Scheme = "http"
			req.URL.Host = koalaHost
			req.Host = koalaHost
			log4go.Debug(req.URL.Scheme)
			log4go.Debug(req.URL)
			log4go.Debug(req)
		}
	*/

	remote, err := url.Parse(baseUrl)
	c.Request.Host = baseUrl[7:]
	//remote, err := url.Parse("http://192.168.18.2:8080")
	if err != nil {
		log4go.Error(err)
		c.Writer.Write([]byte(err.Error()))
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	//buffer, _ := ioutil.ReadAll(c.Request.Body)
	//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
	//proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)

	return nil
}
//var baseUrl = ""
//var jar, _ = cookiejar.New(nil)

func ReverseProxy(c *gin.Context) error {

	//target := baseUrl + c.Request.RequestURI

	/*
		director := func(req *http.Request) {
			r := c.Request
			req = r
			log4go.Debug(req.URL.Scheme)
			req.URL.Scheme = "http"
			req.URL.Host = koalaHost
			req.Host = koalaHost
			log4go.Debug(req.URL.Scheme)
			log4go.Debug(req.URL)
			log4go.Debug(req)
		}
	*/

	remote, err := url.Parse(baseUrl)
	//remote, err := url.Parse("http://192.168.18.2:8080")
	if err != nil {
		log4go.Error(err)
		c.Writer.Write([]byte(err.Error()))
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	//buffer, _ := ioutil.ReadAll(c.Request.Body)
	//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
	//proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)

	return nil
}

func ReverseProxy3(c *gin.Context) error {

	//target := baseUrl + c.Request.RequestURI

	/*
		director := func(req *http.Request) {
			r := c.Request
			req = r
			log4go.Debug(req.URL.Scheme)
			req.URL.Scheme = "http"
			req.URL.Host = koalaHost
			req.Host = koalaHost
			log4go.Debug(req.URL.Scheme)
			log4go.Debug(req.URL)
			log4go.Debug(req)
		}
	*/
	Url := "http://hz91zo.oicp.vip:9980"
	remote, err := url.Parse(Url)

	//c.Request.Host = "127.0.0.1:5600"
	c.Request.Host = "hz91zo.oicp.vip:9980"
	//remote, err := url.Parse("http://192.168.18.2:8080")
	if err != nil {
		log4go.Error(err)
		c.Writer.Write([]byte(err.Error()))
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	//buffer, _ := ioutil.ReadAll(c.Request.Body)
	//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
	//proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)

	return nil
}
func GetScreenList(c *gin.Context)error {
	return ReverseProxy(c)
}

func SetDisplayConfig(c *gin.Context) error {
	return ReverseProxy(c)
}

/*
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
*/
