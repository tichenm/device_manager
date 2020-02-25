package model

// camera 基础表字段
type Camera struct {
	ID   int  `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"id"`
	Camera_address string	`gorm:"column:camera_address" json:"camera_address"`
	Camera_status int	`gorm:"column:camera_status" json:"camera_status"`
	Camera_position string	`gorm:"column:camera_position" json:"camera_position"`
	Camera_type string	`gorm:"column:camera_type" json:"camera_type"`
}

// account 基础表字段
type Account struct {
	ID   int  `gorm:"primary_key:AUTO_INCREMENT;column:id;not null"`
	Ip_address string	`gorm:"column:ip_address"`
	Account string	`gorm:"column:account"`
	Password string	`gorm:"column:password"`
	Activation int	`gorm:"column:activation"`
}