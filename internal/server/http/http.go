package http

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"strings"
	"time"
	"zhiyuan/scaffold/configs"
	"zhiyuan/scaffold/internal/koala"
	"zhiyuan/scaffold/internal/model"
	"zhiyuan/scaffold/service"
)

var (
	svc *service.Service
)
var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
//var (
//	svc *service.Service
//)

// New new a gin server.
func New() {
	var (
		hc struct {
			Server *model.ServerConfig
		}
		cpath string
		//screen struct{
		//	ac  *paladin.Map
		//}
	)

	// 初始化
	if runtime.GOOS == "linux" {
		cpath = "./configs/http.toml"
		if runtime.GOARCH == "arm" {
			cpath = "./configs/http.toml"
		} else if runtime.GOARCH == "amd64" {
			cpath = "./configs/http.toml"
		}
	} else if runtime.GOOS == "windows" {
		cpath = "F:/program/code/go/go_project/src/zhiyuan/device_server/raying_api/configs/http.toml"
	}

	if _, err := toml.DecodeFile(cpath, &hc); err != nil {
		log.Error("read toml file error(%v)", err)
	}

	svc = service.New(configs.Conf)
	//svc
	svc.Createserver()
	account, _ := svc.GetAccount()
	koala.Init(account.Ip_address)
	username := "admin@91zo.com"
	password := "123456"
	koala.KoalaLogin(username, password)
	go Keepalive(username, password)
	go manager.start()
	engine := gin.Default()
	initRouter(engine)
	gin.SetMode(gin.ReleaseMode)
	engine.Run(hc.Server.Addr)
}

func Keepalive(username, password string) {

	for {
		time.Sleep(180 * time.Minute)
		koala.KoalaLogin(username, password)
	}
}

//func Keepalive(username,password string) {
//	//time.Sleep(1800*time.Second)
//	cronTarget := cron.New()
//	spec := "* */30 * * * ?"
//	cronTarget.AddFunc(spec, func() {
//		koala.KoalaLogin(username,password)
//	})
//	cronTarget.Start()
//	log.WithFields(log.Fields{
//		"koala-cron": "checkalive",
//	}).Info("checkalive cron start")
//}
func initRouter(e *gin.Engine) {

	e.Use(Cors())
	system := e.Group("/v1")
	{
		system.GET("/start", howToStart)
		system.POST("/login", Login)
		system.GET("/status", ReverseProxy_koalamate)
		system.GET("/status/check_network", ReverseProxy_koalamate)
		system.PUT("/tatus/sys_time", ReverseProxy_koalamate)
		system.PUT("/status/reboot", ReverseProxy_koalamate)
		system.GET("/config/ip", ReverseProxy_koalamate)
		system.PUT("/config/ip", ReverseProxy_koalamate)
		system.GET("/logs", ReverseProxy_koalamate)
		system.GET("/log/", ReverseProxy_koalamate)
		system.POST("/update", ReverseProxy_koalamate)

	}
	e.GET("/video", wsEndpoint)
	e.GET("/support", wsSupport)
	//e.POST("/callback",wsData)
	service := e.Group("/system")
	{
		service.POST("/screen", CreateCamera)
		service.PUT("/screen/:id", UpdateCamera)
		service.DELETE("/screen/:id", DeleteCamera)
		service.POST("/screens", GetCameras)

	}
	account := e.Group("/account")
	{
		account.GET("/accounts", GetAccount)
		account.PUT("/account/:id", UpdateAccount)

	}
	screen := e.Group("/screen")
	{
		// 获取所有可以相机列表
		screen.GET("/get-screen-list", Get_screen_list)
		// 设置显示设备需要弹窗和显示的相机
		//screen.POST("/set-display-config", ReverseProxy)
		screen.POST("/set-display-config", Set_display_config)
		screen.POST("/set-device-info", ReverseProxy)
		// 获取最新的显示设备设置/screen/set-display-config
		//screen.GET("/get-display-config", ReverseProxy)

		screen.GET("/get-display-config", Get_display_config)

		screen.GET("/get", ReverseProxy)
		screen.POST("/set-theme", ReverseProxy)
		screen.GET("/get_theme_config", ReverseProxy)
		screen.GET("/theme", ReverseProxy)
		screen.POST("/add-visitor", ReverseProxy)
		screen.GET("/weather", ReverseProxy)
		screen.GET("/avatars", ReverseProxy)
		screen.GET("/custom.css", ReverseProxy)
		screen.GET("/vip-cards.css", ReverseProxy)
		screen.GET("/custom.html", ReverseProxy)
	}
	koala_static := e.Group("/koala_static")
	{
		koala_static.Any("/*path", KoalaStatic)
	}
	e.GET("/getServerTime", GetServerTime)

}

// example for http request handler.
func howToStart(c *gin.Context) {
	c.String(0, "Golang 大法好 !!!")
}

func Login(c *gin.Context) {
	code := 0
	err_msg := ""

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		code = -100
		err_msg = "用户名或密码不能为空"
	}

	if username != "admin" || password != "zybox" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"err_msg": "用户名或密码错误",
		})
		return
	}

	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"err_msg": err_msg,
		})
		return
	}

	uid, _ := uuid.NewV4()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    uid.String(),
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
	})
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"err_msg": err_msg,
		"data":    "验证成功",
	})
	return

}

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan *Message_obj
	register   chan *Client
	unregister chan *Client
}

type Message_obj struct {
	id   string
	data []byte
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
	dutie  int
}

var manager = ClientManager{
	broadcast:  make(chan *Message_obj),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func wsEndpoint(c *gin.Context) {

	url_RPC := c.Query("url")
	ip := svc.FindUrl(url_RPC)
	port_id := svc.GeneratePort(ip)
	log.Println(port_id)
	w := c.Writer
	r := c.Request

	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Println("Client Successfully Connected...")

	go activateevent(port_id)
	client := &Client{id: port_id, socket: ws, send: make(chan []byte), dutie: 0}

	manager.register <- client
	go client.write()

}

func activateevent(port string) {
	callback_url := "http://127.0.0.1:" + port + "/cilent"
	req, err := http.NewRequest("GET", callback_url, nil)
	if err != nil {
		//log4go.Error(err)
		return
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		//log4go.Error(err.Error())
	}
}

func wsSupport(c *gin.Context) {

	port := c.Query("port")
	var flag bool
	log.Println(port)
	w := c.Writer
	r := c.Request
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Println("Client Successfully Connected Support")
	client_support := &Client{id: port, socket: ws, send: make(chan []byte), dutie: 1}
	if len(manager.clients) == 0 {
		manager.register <- client_support
		go client_support.read()
	} else {
		for conn := range manager.clients {
			if conn.id == port && conn.dutie == 1 {
				flag = true
			}
		}
		if flag {
			ws.Close()
			return
		} else {
			manager.register <- client_support
			go client_support.read()
		}
	}

}

func (c *Client) write() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("写关闭了")
	}()

	for {
		select {
		case message, ok := <-c.send: //这个管道有了数据 写这个消息出去
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("发送关闭提示")
				return
			}

			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				manager.unregister <- c
				c.socket.Close()
				fmt.Println("写不成功数据就关闭了")
				break
			}
			fmt.Println("写数据")
		}
	}
}
func (c *Client) read() {

	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("读关闭")
	}()

	for {

		_, message, err := c.socket.ReadMessage()
		log.Println(string(message))
		fmt.Println("是在不停的读吗？")
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			fmt.Println("读不到数据就关闭？")
			break
		}
		jsonMessage, err := json.Marshal(string(message))

		log.Println(jsonMessage)
		log.Println(len(jsonMessage))
		message_obj := &Message_obj{
			id:   c.id,
			data: jsonMessage,
		}
		manager.broadcast <- message_obj //激活start 程序 入广播管道
		fmt.Println("发送数据到广播")

	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: //新客户端加入
			manager.clients[conn] = true
			log.Println(" a new socket has connected.")

		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				log.Println("a socket has disconnected.")

			}
		case message := <-manager.broadcast: //读到广播管道数据后的处理
			fmt.Println(string(message.data))
			for conn := range manager.clients {
				fmt.Println("当前客户端", conn.id)
				if conn.dutie == 0 && conn.id == message.id {
					select {
					case conn.send <- message.data: //调用发送给指定客户端
					default:
						fmt.Println("要关闭连接啊")
						close(conn.send)
						delete(manager.clients, conn)
					}
				}

			}
		}
	}
}

func reader(conn *websocket.Conn) {
	var (
		data []byte
	)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		//获取识别记录FUNCTION
		time.Sleep(500 * time.Millisecond)
		data = []byte(string(p))

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println(err)
			return
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
