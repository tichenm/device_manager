package model

import (
	xtime "github.com/bilibili/kratos/pkg/time"
)

// Kratos hello kratos.
type KoalaPerson_ struct {
	Id               int    `json:"id"`
	Company_id       int    `json:"company_id"`
	Create_time      int64  `json:"create_time"`
	InterView_pinyin string `json:"interView_pinyin"`
	Visitor_type     string `json:"visitor_type"`
	Subject_type     int    `json:"subject_type"`
	Email            string `json:"email"`
	Password_reseted string `json:"password_reseted"`
	Name             string `json:"name"`
	Pinyin           string `json:"pinyin"`
	Gender           int    `json:"gender"`
	Photo_ids        string `json:"photo_ids"`
	Phone            string `json:"phone"`
	Avatar           string `json:"avatar"`
	Department       string `json:"department"`  //部门
	Title            string `json:"title"`       //职位
	Description      string `json:"description"` //签名
	Job_number       string `json:"job_number"`  //工号
	Remark           string `json:"remark"`
	Birthday         int64  `json:"birthday"`     //生日	时间戳（秒）
	Entry_date       int64  `json:"entry_date"`   //入职时间	时间戳（秒）
	Purpose          int64  `json:"purpose"`      //(访客属性) 来访目的	0: 其他, 1: 面试, 2: 商务, 3: 亲友, 4: 快递送货
	Interviewee      int64  `json:"interviewee"`  //(访客属性) 受访人
	Come_from        int64  `json:"come_from"`    //(访客属性) 来访单位
	Start_time       int64  `json:"start_time"`   //(访客属性) 预定来访时间	时间戳（秒）
	End_time         int64  `json:"end_time"`     //(访客属性) 预定离开时间	时间戳（秒）
	Visit_notify     string `json:"visit_notify"` //(访客属性) 来访是否发APP消息推送
}

type Koala_person struct {
	Avatar             string                `json:"avatar"`
	Birthday           int64                 `json:"birthday"`
	Come_from          string                `json:"come_from"`
	Company_id         int                   `json:"company_id"`
	Create_time        int64                 `json:"create_time"`
	Department         string                `json:"department"`
	Description        string                `json:"description"`
	Email              string                `json:"email"`
	End_time           int64                 `json:"end_time"`
	Entry_date         int64                 `json:"entry_date"`
	Gender             int                   `json:"gender"`
	Id                 int                   `json:"id"`
	Interviewee        string                `json:"interviewee"`
	Interviewee_pinyin string                `json:"interviewee_pinyin"`
	Job_number         string                `json:"job_number"`
	Name               string                `json:"name"`
	Password_reseted   bool                  `json:"password_reseted"`
	Phone              string                `json:"phone"`
	Photos             []Koala_person_photos `json:"photos"`
	Pinyin             string                `json:"pinyin"`
	Purpose            int                   `json:"purpose"`
	Remark             string                `json:"remark"`
	Start_time         int64                 `json:"start_time"`
	Subject_type       int                   `json:"subject_type"`
	Title              string                `json:"title"`
	Visit_notify       bool                  `json:"password_reseted"`
}
type Koala_person_photos struct {
	Origin_url string  `json:"origin_url"`
	Company_id int     `json:"company_id"`
	Id         int     `json:"id"`
	Quality    float32 `json:"quality"`
	Subject_id int     `json:"subject_id"`
	Url        string  `json:"url"`
	Version    int     `json:"version"`
}

type ServerConfig struct {
	Network      string         `dsn:"network"`
	Addr         string         `dsn:"address"`
	Timeout      xtime.Duration `dsn:"query.timeout"`
	ReadTimeout  xtime.Duration `dsn:"query.readTimeout"`
	WriteTimeout xtime.Duration `dsn:"query.writeTimeout"`
}

type FaceInfo struct {
	ObjectID    int32
	BoundingBox []int
	PosePitch   float32
	PoseRoll    float32
	Blur        float32
	Sex         string
	Age         int
	Minority    int
}

type FaceRecognizeInfo struct {
	FaceToken       string
	SearchScore     float32
	SearchThreshold float32
	PersonInfo      *FaceInfo
}

type FaceRecognizeGroup struct {
	GroupAlias string
}

type PersonInfo struct {
	Name            string
	Birthday        string
	Sex             string
	CertificateType string
	ID              string
	Country         string
	Province        string
	City            string
}

type EventCommInfo struct {
	Resolution     []int
	PictureType    int
	MachineAddress string
	SerialNo       string
}

type IDCardInfo struct {
	Name           string
	Sex            string
	Nation         string
	Number         string
	Address        string
	Office         string
	ValidTimeStart string
	ValidTimeStop  string
	ProfilePic     string
}

//type Addperson struct {
//	GroupID	int	`json:"GroupID"`
//	PersonInfo	interface{}	`json:"PersonInfo"`
//	ImageInfo interface{}	`json:"ImageInfo"`
//}

type Addperson struct {
	//GroupID	int	`json:"GroupID"`
	PersonInfo interface{} `json:"Person"`
	//ImageInfo interface{}	`json:"ImageInfo"`
}
type CreatePhoto struct {
	GroupID    int         `json:"GroupID"`
	PersonInfo interface{} `json:"PersonInfo"`
	ImageInfo  interface{} `json:"ImageInfo"`
}

type Rpcmodel struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Rpcmodel_findperson struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Object int         `json:"object"`
	Params interface{} `json:"params"`
}

type NormalResponsemodel struct {
	Code    int         `json:"code"`
	Message string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

type RYPerson struct {
	Code int    `json:"Code"`
	Name string `json:"Name"`
	Sex  string `json:"Sex"`
	Type int    `json:"Type"`
	//Country				string	`json:"Country"`
	//Province					string	`json:"Province"`
	//City					string	`json:"City"`
	CertificateType string `json:"CertificateType"`
	GroupName       string `json:"GroupName"`
	Birthday        string `json:"Birthday"`
}
type RYPerson_2 struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Sex  string `json:"Sex"`
	//Type	int	`json:"Type"`
	Country         string `json:"Country"`
	Province        string `json:"Province"`
	City            string `json:"City"`
	CertificateType string `json:"CertificateType"`
	//GroupName				string	`json:"GroupName"`
	//Birthday	string	`json:"Birthday"`
}
type RYimg struct {
	Lengths [1]int64 `json:"Lengths"`
	Amount  int      `json:"Amount"`
}

type Pair struct {
	a, b, c interface{}
}

//统一的请求
type Person struct {
	Subject_type int    `json:"subject_type"`
	Subject_Id   int    `json:"subject_id"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Ip           string `json:"ip"`
}

//
type Resp4Device struct {
	Code    int         `json:"code"`
	Err_msg string      `json:"err_msg"`
	Data    interface{} `json:"data"`
	Page    interface{} `json:"Page"`
}

type EventPerson struct {
	PersonName string	`json:"personName"`
	SubjectId string	`json:"subjectId"`
	PersonType string	`json:"personType"`
	FacePicture string	`json:"facePicture"`
	Similarity string	`json:"similarity"`
	CameraIp string		`json:"cameraIp"`
	Timestamp int64		`json:"timestamp"`
}