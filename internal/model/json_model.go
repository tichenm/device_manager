package model

// camera 基础表字段
type Camera_json struct {
	Camera_address string	`json:"camera_address"`
	Camera_position string	`json:"camera_position"`
	Camera_type string	`json:"camera_type"`
}

// account 基础表字段
type Account_josn struct {

	Ip_address string	`json:"ip_address"`
	Account string	`json:"account"`
	Password string	`json:"password"`
}

// camera 基础表字段
type Camera_status_json struct {
	Camera_status int	`json:"camera_status"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}