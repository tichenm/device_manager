package model

// camera 基础表字段
type Camera struct {
	ID   int  `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"id"`
	Camera_address string	`gorm:"column:camera_address" json:"camera_address"`
	Camera_status int	`gorm:"column:camera_status" json:"camera_status"`
	Camera_position string	`gorm:"column:camera_position" json:"camera_position"`
	Camera_type string	`gorm:"column:camera_type" json:"camera_type"`
	Camera_RTSP string  `gorm:"column:camera_rtsp" json:"camera_rtsp"`
	Camera_token string  `gorm:"column:camera_token" json:"camera_token"`
}

// account 基础表字段
type Account struct {
	ID   int  `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"id"`
	Ip_address string	`gorm:"column:ip_address" json:"ip_address"`
	Account string	`gorm:"column:account" json:"account"`
	Password string	`gorm:"column:password" json:"password"`
	Activation int	`gorm:"column:activation" json:"activation"`
}