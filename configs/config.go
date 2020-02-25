package configs

import (
	//"go-common/library/net/rpc"
	"go-common/library/database/sql"
	"go-common/library/log"
	//bm"github.com/bilibili/kratos/pkg/net/http/blademaster"
	//"github.com/bilibili/kratos/database/orm"
	//"go-common/library/database/orm"
	//"go-common/library/database/orm"
	"go-common/library/database/orm"
)

var (
	// ConfPath local config path
	// Conf config
	Conf   = &Config{}
	confPath string
)

// Config str
type Config struct {
	// base
	// channal len
	Mail          *Mail
	// log
	Log *log.Config
	// identify
	// tracer
	// tick load pgc
	//Tick time.Duration
	// orm
	ORM *orm.Config
	// host
	//Host *Host
	//// Bfs
	//Bfs *Bfs
	// http client
	//HTTPClient *bm.ClientConfig
	//// image client
	//ImageClient *bm.ClientConfig
	// BM HTTPServers
	//Mysql *sql.Config
	//BM *bm.ServerConfig
	// db
	// budget
	//Budget *Budget
	//// rpc client
	//VipRPC *rpc.ClientConfig
	//// grpc client
	//Account *warden.ClientConfig
	//// shell config
	//ShellConf *ShellConfig
	//OtherConf *OtherConfig
}

//DB .
type breaker struct {
	Archive     *sql.Config
	ArchiveRead *sql.Config
	Manager     *sql.Config
	Oversea     *orm.Config
	Creative    *sql.Config
}


//Mail ...
type Mail struct {
	Host     string
	Port     int
	From     string
	Password string
	To       []string
}
