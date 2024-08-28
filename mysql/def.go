package mysql

import (
	libSQL "database/sql"
	"fmt"
	"time"

	// init driver for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/tonny-zhang/goutils/logger"
)

const slowTime = 200 * time.Millisecond

// Config mysql 的配置
type Config struct {
	MysqlUser string
	MysqlPwd  string
	MysqlHost string
	MysqlPort int
	MysqlDB   string

	MaxOpenConns int

	LogPrefix string
}

// Conf 初始化配置
func Conf(conf Config) (db *DB, e error) {
	logPrefix := conf.LogPrefix
	if logPrefix == "" {
		logPrefix = "[mysql]"
	}
	log := logger.PrefixLogger(logPrefix)
	var user, pwd, host, database string
	var port int
	host = conf.MysqlHost
	port = conf.MysqlPort
	user = conf.MysqlUser
	pwd = conf.MysqlPwd
	database = conf.MysqlDB

	dbFlag := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, pwd, host, port, database)
	conn, e := libSQL.Open("mysql", dbFlag)
	if e == nil && conf.MaxOpenConns > 0 {
		conn.SetMaxOpenConns(conf.MaxOpenConns)
		conn.SetMaxIdleConns(conf.MaxOpenConns)
	}

	// log.Info("mysql初始化 %s:%v", host, port)

	if e == nil {
		db = &DB{
			db:  conn,
			Log: log,
		}
	}
	return
}
