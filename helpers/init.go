package helpers

import (
	"encoding/json"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
)

var MyDB *sqlx.DB

type config struct {
	App struct {
		Port string `ini:"port"`
	} `ini:"app"`
	Mysql struct {
		Dsn string `ini:"dsn"`
	} `ini:"mysql"`
	Auth struct {
		SecretKey string `ini:"secretKey"`
		Username  string `ini:"username"`
		Password  string `ini:"password"`
	} `ini:"auth"`
}

func (data config) GetAsJSON() (res string) {
	parsed, _ := json.Marshal(data)
	res = string(parsed)
	return
}

var Config = new(config)

var log = logging.MustGetLogger("mail")

func ConfigInit() {

	ini.MapTo(&Config, "config.ini")

	log.Debug("Loaded config: ", Config.GetAsJSON())
}

func InitMysql() {
	log.Debug("Connect to MySQL", Config.Mysql.Dsn)
	MyDB = sqlx.MustConnect("mysql", Config.Mysql.Dsn)
	MyDB.SetMaxIdleConns(10)
	MyDB.SetMaxOpenConns(20)
}
