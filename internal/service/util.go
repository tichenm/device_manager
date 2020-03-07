package service

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"time"
	log "github.com/sirupsen/logrus"
	"zhiyuan/device_server/raying_api/internal/model"
)

// Retry function. fn is the retry function.
func Retry(attempts int, sleep time.Duration, fn func() bool) bool {
	if err := fn(); err != true {

		if attempts--; attempts > 0 {
			log.Warnf("retry func error: %s. attemps #%d after %s.", "failed", attempts, sleep)
			time.Sleep(60 * time.Second)
			return Retry(attempts, 2*sleep, fn)
		}

		return err
	}
	//IsConnected = true
	return true
}

func Buildjson(params interface{}) (*simplejson.Json, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.WithFields(log.Fields{
			"raying": "Marshal_json_err",
		}).Info(err)
		return nil, err
	}
	jsondata, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"raying": "simplejson_NewJson_err",
		}).Info(err)
		return nil, err
	}
	return jsondata, err

}

func Json2Map(params interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.WithFields(log.Fields{
			"raying": "Marshal_json_err",
		}).Info(err)
		return nil, err
	}
	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal(jsonBytes, &mapResult); err != nil {
		log.WithFields(log.Fields{
			"raying": "simplejson_NewJson_err",
		}).Info(err)
		return nil, err
	}
	return mapResult, err

}

func Datatransform(code int, msg string, data interface{}) (*simplejson.Json) {
	var (
		resp4Device model.Resp4Device
	)
	resp4Device.Code = code
	resp4Device.Err_msg = msg
	resp4Device.Data = data
	jsondata, _ := Buildjson(resp4Device)
	return jsondata
}
