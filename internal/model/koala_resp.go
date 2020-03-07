package model



type Screens struct {
	Box_address string	`json:"box_address"`
	Box_heartbeat int64	`json:"box_heartbeat"`
	Box_id int	`json:"box_id"`
	Box_status string	`json:"box_status"`
	Box_token	string	`json:"box_token"`
	Camera_address string	`json:"camera_address"`
	Camera_name string	`json:"camera_name"`
	Camera_position string	`json:"camera_position"`
	Camera_status string	`json:"camera_status"`
	Description string	`json:"description"`
	Id int	`json:"id"`
	Is_select int	`json:"is_select"`
	Network_switcher string	`json:"network_switcher"`
	Network_switcher_drive int	`json:"network_switcher_drive"`
	Network_switcher_status string	`json:"network_switcher_status"`
	Network_switcher_token string	`json:"network_switcher_token"`
	Screen_token string	`json:"screen_token"`
	Type int	`json:"type"`
}

type Display_device struct {
	Background	string `josn:"background"`
	Background_layout	int `josn:"background_layout"`
	Card_duration	int `josn:"card_duration"`
	Card_theme	string `josn:"card_theme"`
	Card_theme_vip	string `josn:"card_theme_vip"`
	City	string `josn:"city"`
	Description	string `josn:"description"`
	Display_info_timestamp	int `josn:"display_info_timestamp"`
	Heartbeat	int64 `josn:"heartbeat"`
	Id	int `josn:"id"`
	Logo	string `josn:"logo"`
	Reload_timestamp	int `josn:"reload_timestamp"`
	Screens	[]Screens `josn:"screens"`
	Status	string `josn:"status"`
	Stranger_warn	int `josn:"stranger_warn"`
	Theme	string `josn:"theme"`
	Token	string `josn:"token"`
	User_info_timestamp	int `josn:"user_info_timestamp"`
	Video_screen_id	int `josn:"video_screen_id"`
	Yellowlist_warn	int `josn:"yellowlist_warn"`
}

