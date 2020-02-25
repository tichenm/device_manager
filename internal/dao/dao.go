package dao

//import "gopkg.in/gomail.v2"

//import (
//	//"github.com/jinzhu/gorm"
//	xsql "go-common/library/database/sql"
//	"zhiyuan/ai_dormitory_apis/school_affairs/conf"
//	"gopkg.in/gomail.v2"
//	bm "go-common/library/net/http/blademaster"
//	"github.com/bilibili/kratos/pkg/database/sql"
//)
//// Dao dao
//type Dao struct {
//	c             *conf.Config
//	//redis         *redis.Pool
//	//bfredis       *redis.Pool
//	db            *xsql.DB
//	dbCms         *xsql.DB
//	HTTPClient    *bm.Client
//	//SearchClient  searchv1.SearchClient
//	//VideoClient   videov1.VideoClient
//	//AccountClient account.AccountClient
//	email         *gomail.Dialer
//	//noticeClient  notice.NoticeClient
//
//}
//
//package dao

import (
	"gopkg.in/gomail.v2"
	"zhiyuan/scaffold/configs"
	"zhiyuan/scaffold/internal/model"
	//"go-common/library/database/sql"
	//	"go-common/library/database/sql"
	//bm"go-common/library/net/http/blademaster"
	"context"
	"github.com/alecthomas/log4go"
	"go-common/library/log"
	//"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm"
	"go-common/library/database/sql"
	"runtime"
	"strconv"
	"zhiyuan/ai_dormitory_apis/config"
)





type Dao struct {
	c             *configs.Config
	//redis         *redis.Pool
	//bfredis       *redis.Pool
	//creativeDB       *sql.DB
	db *sql.DB
	// orm
	crmdb *gorm.DB
	//rddb *sql.DB
	//dbCms         *xsql.DB
	//HTTPClient    *bm.Client
	//SearchClient  searchv1.SearchClient
	//VideoClient   videov1.VideoClient
	//AccountClient account.AccountClient
	email         *gomail.Dialer
	//noticeClient  notice.NoticeClient
}


//// New init mysql db
//func New(c *conf.Config) (dao *Dao) {
//	dao = &Dao{
//		c:             c,
//		//redis:         redis.NewPool(c.Redis),
//		//bfredis:       redis.NewPool(c.BfRedis),
//		//db:            xsql.NewMySQL(c.MySQL),
//		//dbCms:         xsql.NewMySQL(c.MySQLCms),
//		//HTTPClient:    bm.NewClient(c.BM.Client),
//		//SearchClient:  newSearchClient(c.GRPCClient["search"]),
//		//VideoClient:   newVideoClient(c.GRPCClient["video"]),
//		//AccountClient: newAccountClient(c.GRPCClient["account"]),
//		email:         gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.From, c.Mail.Password),
//		//noticeClient:  newNoticeClient(c.GRPCClient["notice"]),
//	}
//	return
//}
//func checkErr(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
// New new a dao.
//func New(c *conf.Config) (d *Dao) {
//	d = &Dao{
//		c:             c,
//		db:            sql.NewMySQL(c.DB.Archive),
//		//rddb:          sql.NewMySQL(c.DB.ArchiveRead),
//		//redis:         redis.NewPool(c.Redis.Track.Config),
//		//hbase:         hbase.NewClient(&c.HBase.Config),
//		//userCardURL:   c.Host.Account + "/api/member/getCardByMid",
//		//addQAVideoURL: c.Host.Task + "/vt/video/add",
//		//clientW:       bm.NewClient(c.HTTPClient.Write),
//		//clientR:       bm.NewClient(c.HTTPClient.Read),
//		creativeDB:    sql.NewMySQL(c.DB.Creative),
//		email:         gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.From, c.Mail.Password),
//	}
//	return d
//}

func New(c *configs.Config)*Dao {

	//Db.Set("gorm:", "ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE utf8_general_ci")

	//if runtime.GOOS == "linux" {
	//	config.Init("/home/zhiyuan/ai_dormitory_apis_formal/conf.yaml")
	//} else if runtime.GOOS == "windows" {
	//	config.Init("../conf.yaml")
	//}
	if runtime.GOOS == "linux" {
		//config.Init("/home/zhiyuan/ai_dormitory_apis/conf.yaml")
		config.Init("./conf.yaml")
	} else if runtime.GOOS == "windows" {
		config.Init("../conf.yaml")
	}
	//if err := koala_api.Init(config.Gconf.KoalaHost, config.Gconf.KoalaPort); err != nil {
	//	return nil
	//}
	//var (
	//	argName = &orm.Config {
	//		//Network :string(config.Gconf.SchoolServerPort),
	//		//Addr : config.Gconf.Addr,
	//		DSN : config.Gconf.Dsn,
	//		//ReadDSN:config.Gconf.ReadDSN,
	//		Active:config.Gconf.Active,
	//		Idle:config.Gconf.Idle,
	//		//IdleTimeout:config.Gconf.IdleTimeout,
	//		//QueryTimeout:config.Gconf.QueryTimeout,
	//		//ExecTimeout:config.Gconf.ExecTimeout,
	//		//TranTimeout:config.Gconf.TranTimeout,
	//	}
	//)
	//var (
	//	argName = &orm.Config {
	//		//Network :string(config.Gconf.SchoolServerPort),
	//		Addr : config.Gconf.Addr,
	//		DSN : config.Gconf.Dsn,
	//		ReadDSN:config.Gconf.ReadDSN,
	//		Active:config.Gconf.Active,
	//		Idle:config.Gconf.Idle,
	//		IdleTimeout:config.Gconf.IdleTimeout,
	//		QueryTimeout:config.Gconf.QueryTimeout,
	//		ExecTimeout:config.Gconf.ExecTimeout,
	//		TranTimeout:config.Gconf.TranTimeout,
	//	}
	//)
	var d = &Dao{
		c:c,
	}
	crmdb, err := gorm.Open("mysql", config.Gconf.DBUsername + ":" + config.Gconf.DBPassword+
		"@tcp("+ config.Gconf.DBIP+ ":"+ strconv.Itoa(config.Gconf.DBPort)+ ")/"+ config.Gconf.DBName+ "?parseTime=true&charset=utf8")
	if crmdb == nil {
		log.Error("connect to db fail, err=%v", err)
		return nil
	}

	crmdb.SingularTable(true)
	d.crmdb = crmdb
	d.crmdb.LogMode(true)
	//d.crmdb.DB().SetMaxIdleConns(argName.Active)
	d.crmdb.DB().SetMaxIdleConns(20)
	//d.crmdb.DB().SetMaxOpenConns(argName.Active)
	d.crmdb.DB().SetMaxOpenConns(20)
	//db, err := gorm.Open("mysql", config.Gconf.DBUsername + ":" + config.Gconf.DBPassword+
	//	"@tcp("+ config.Gconf.DBIP+ ":"+ strconv.Itoa(config.Gconf.DBPort)+ ")/"+ config.Gconf.DBName+ "?parseTime=true&charset=utf8")
	//if err != nil {
	//	//log.Fatalln(err)
	//}

	//db.SingularTable(true)
	////defer db.Close()
	////空闲连接和最大打开连接
	//db.DB().SetMaxIdleConns(20)
	//db.DB().SetMaxOpenConns(20)
	//
	//if err := db.DB().Ping(); err != nil {
	//	log.Fatalln(err)
	//}
	//Db = db
	CreateTable(crmdb)
	return d
}

func CreateTable(db *gorm.DB) {
	if !db.HasTable(model.Camera{}) {
		db.CreateTable(&model.Camera{})
		log4go.Info("创建学生表成功")
	} else {
		db.AutoMigrate(&model.Camera{})
		log4go.Info("更新学生表成功")
	}
	if !db.HasTable(model.Account{}) {
		db.CreateTable(&model.Account{})
		log4go.Info("创建职工表成功")
	} else {
		db.AutoMigrate(&model.Account{})
		log4go.Info("更新职工表成功")
	}
}
// Close close dao.
// Close close resource.
func (d *Dao) Close() {
	if d.crmdb != nil {
		d.crmdb.Close()
	}
}
//if d.creativeDB != nil {
//	d.creativeDB.Close()
//}
//d.redis.Close()
//}


//// BeginTran begin mysql transaction
//func (d *Dao) BeginTran(c context.Context) (*xsql.Tx, error) {
//	return d.db.Begin(c)
//}

// BeginTran begin transcation.
func (d *Dao) GetDb() *gorm.DB {
	return d.crmdb
}

//// Ping dao ping
//func (d *Dao) Ping(c context.Context) error {
//	// TODO: if you need use mc,redis, please add
//	return d.db.Ping(c)
//}
// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
//StartTx start tx
//func (d *Dao) StartTx(c context.Context) (tx *sql.Tx, err error) {
//	if d.db != nil {
//		tx, err = d.db.Begin(c)
//	}
//	return
//}