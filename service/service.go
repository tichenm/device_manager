package service

//
import(
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	//"zhiyuan/ai_dormitory_apis/school_affairs/conf"
	"zhiyuan/scaffold/internal/dao"
	"context"
	"zhiyuan/scaffold/configs"
	log "github.com/sirupsen/logrus"
)
//// Service service.
type Service struct {
	//	ac  *paladin.Map
	//	dao *dao.Dao
	c           *configs.Config
	dao         *dao.Dao
}
//
// New new a service and return.
func New(c *configs.Config) (s *Service) {
	//var ac = new(paladin.TOML)
	//if err := paladin.Watch("application.toml", ac); err != nil {
	//	panic(err)
	//}
	s = &Service{
		c:  c,
		dao: dao.New(c),
	}
	return s
}
//
// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}
//
// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) Createserver() {
	camera_type := map[string]string{
		//"巨龙":"/zybox/camera_server/JL/JL_server",
		"巨龙":"F:/program/code/go/go_project/src/client_test/main.exe",
		"锐颖":"/zybox/camera_server/RY/RY_server",
		"地平线":"/zybox/camera_server/DPX/DPX_server",
	}

	result,_,_,_ := s.dao.GetCameras("",0,1,1000)
	//account_obj := s.dao.GetAccount
	for i:=0;i<len(result);i++{
		//PORT , FLAG := s.ReadPorttoml()
		//if FLAG{
		//	PORT += 1
		//}else{
		//	PORT=4089
		//}
		PORT := s.GeneratePort(result[i].Camera_address)
		//s.SavePorttoml(PORT)
		url := camera_type[result[i].Camera_type]
		//account := account_obj.Account
		//password :=account_obj.password
		//address := account_obj.address
		account := "admin@91zo.com"
		password := "123456"
		address := "192.168.18.50"
		port :=  PORT
		cmd := exec.Command(url,account,password,address,port)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Print(string(stdout))
	}
}

func (s *Service)ReadPorttoml() (int,bool) {

	type Port_obj struct {
		Port int
	}
	var port Port_obj
	path_GI :="./port.toml"

	if _, err := toml.DecodeFile(path_GI, &port); err != nil {
		log.Error("read toml file error(%v)", err)
		return  0,false
	}

	return  port.Port,true
}

func (s *Service)SavePorttoml( port_value int) (bool) {
	type Port_obj struct {
		Port int
	}
	var port Port_obj
	port.Port = port_value
	file_path :="./port.toml"
	f, err := os.Create(file_path)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if err := toml.NewEncoder(f).Encode(port); err != nil {
		// failed to encode
		log.Fatal(err)
		return false
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		log.Fatal(err)
		return false
	}
	return true
}

func (s *Service)GeneratePort( address string) (string) {
	starter := address
	numberarr := strings.Split(starter, ".")
	variable,_ := strconv.Atoi(numberarr[3])
	sort.Strings(numberarr)
	larger,_ := strconv.Atoi(numberarr[3])
	PORT := variable*larger
	if PORT < 10000 {
		PORT += 10000
	}
	return strconv.Itoa(PORT)
}


func (s *Service)FindUrl(str string) string {
	// 创建一个正则表达式匹配规则对象
	if str == ""{
		return  ""
	}

	regular :=`(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`
	reg := regexp.MustCompile(regular)
	// 利用正则表达式匹配规则对象匹配指定字符串
	res := reg.FindAllString(str, -1)
	if(res == nil){
		return  "3.14"
	}
	ip := strings.Join(res,".")
	return  ip
}




//func main() {
//	app := "echo"
//
//	arg0 := "-e"
//	arg1 := "Hello world"
//	arg2 := "\n\tfrom"
//	arg3 := "golang"
//
//	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
//	stdout, err := cmd.Output()
//
//	if err != nil {
//		Println(err.Error())
//		return
//	}
//
//	Print(string(stdout))
//}


