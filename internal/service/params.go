package service

type Sessionmaker struct {
	Cookies       string
	Authorization string
}
type Facegroupuser struct {
	GroupID int
}

var Result struct {
	Cookies       string
	Authorization string
	GroupID       int
}
//var Obj map[string]interface{}
var Callback_url = ""
//var Authorization = ""
//var GroupID = 1
