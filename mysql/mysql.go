package mysql

import (
	"context"
	libSQL "database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"time"

	// init driver for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/tonny-zhang/goutils/logger"
)

var db *libSQL.DB
var log logger.Logger

const slowTime = 200 * time.Millisecond

// ErrorNotExists 不存在
type ErrorNotExists struct {
	errMsg string
	error
}

func (e ErrorNotExists) Error() string {
	return e.errMsg
}

// Config mysql 的配置
type Config struct {
	MysqlUser string
	MysqlPwd  string
	MysqlHost string
	MysqlPort int
	MysqlDB   string

	MaxOpenConns int
}

func init() {
	var user, pwd, host, port, database string
	log = logger.PrefixLogger("[mysql]")
	if db == nil {
		if v := os.Getenv("MYSQL_HOST"); v != "" {
			host = v
		}
		if v := os.Getenv("MYSQL_PORT"); v != "" {
			port = v
		}
		if v := os.Getenv("MYSQL_USER"); v != "" {
			user = v
		}
		if v := os.Getenv("MYSQL_PWD"); v != "" {
			pwd = v
		}
		if v := os.Getenv("MYSQL_DATABASE"); v != "" {
			database = v
		}
		var err error
		dbFlag := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, database)
		db, err = libSQL.Open("mysql", dbFlag)
		if err != nil {
			fmt.Println(err)
		} else {
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)
		}
	}
	log.Info("mysql初始化 %s:%v", host, port)

}

// Conf 初始化配置
func Conf(conf Config) {
	log = logger.PrefixLogger("[mysql]")
	var user, pwd, host, database string
	var port int
	if db == nil {
		host = conf.MysqlHost
		port = conf.MysqlPort
		user = conf.MysqlUser
		pwd = conf.MysqlPwd
		database = conf.MysqlDB

		var err error
		dbFlag := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, pwd, host, port, database)
		db, err = libSQL.Open("mysql", dbFlag)
		if err != nil {
			fmt.Println(err)
		} else {
			db.SetMaxOpenConns(conf.MaxOpenConns)
			db.SetMaxIdleConns(conf.MaxOpenConns)
		}
	}
	log.Info("mysql初始化 %s:%v", host, port)
}

// 记录慢sql
func traceSlowSQL(sql string, cost time.Duration) {
	if cost > slowTime {
		fileNumInfo := ""
		_, file, line, _ := runtime.Caller(3)
		fileNumInfo = fmt.Sprintf("%s:%d", file, line)
		log.Warn("%s SLOW SQL %v>= %v [%s]", fileNumInfo, cost, slowTime, sql)
	}
}

func list(sql string, args ...any) (rows *libSQL.Rows, query string, err error) {
	if db != nil {
		start := time.Now()

		sql, err = interpolateParams(sql, args...)
		if err == nil {
			rows, err = db.Query(sql)
			// rows, err = db.Query(sql, args...)
			query = sql
		}
		cost := time.Since(start)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		if errmsg != "" {
			log.Error("%-12s\t[%s] err: %s", cost, sql, errmsg)
		} else {
			traceSlowSQL(sql, cost)
			log.Info("%-12s\t[%s]", cost, sql)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}

	return
}

var nilpointer any

var tagMapCache = make(map[string]map[string]int)

// https://studygolang.com/topics/9280
func strutForScan(columns []string, u any) []any {
	val := reflect.ValueOf(u).Elem()

	tagMap := make(map[string]any)
	typeReflect := reflect.TypeOf(u).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := typeReflect.Field(i)
		key := typeField.Tag.Get("db")
		if key != "" {
			tagMap[key] = valueField.Addr().Interface()
		}
	}

	valuePointers := make([]any, len(columns))
	for i, c := range columns {
		v, ok := tagMap[c]
		if ok {
			valuePointers[i] = v
		} else {
			valuePointers[i] = &nilpointer
		}

	}
	return valuePointers
}

// NOTICE: 这里加上PkgPath之后性能和原来差不多，如果不加就有struct名冲突的危险
//
//	后续如果性能差太多就得用原生去处理
func strutForScanCache(columns []string, u any) []any {
	val := reflect.ValueOf(u).Elem()

	// fmt.Println(typeReflect.Name(), typeReflect.PkgPath())
	// TODO: 是否能得到包在内的全称
	t := val.Type()
	key := t.PkgPath() + t.Name()
	// key := val.Type().Name()

	// fmt.Println(key, val.Type().PkgPath())
	keys, ok := tagMapCache[key]
	tagMap := make(map[string]any)
	if !ok {
		keys := make(map[string]int, 0)
		// 对同类型进行缓存
		typeReflect := reflect.TypeOf(u).Elem()
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			typeField := typeReflect.Field(i)
			key := typeField.Tag.Get("db")
			if key != "" {
				tagMap[key] = valueField.Addr().Interface()
				keys[key] = i
			}
		}
		tagMapCache[key] = keys
	} else {
		for k, index := range keys {
			tagMap[k] = val.Field(index).Addr().Interface()
		}
	}

	valuePointers := make([]any, len(columns))
	for i, c := range columns {
		v, ok := tagMap[c]
		if ok {
			valuePointers[i] = v
		} else {
			valuePointers[i] = &nilpointer
		}

	}
	return valuePointers
}

// ChangeLog 更改日志
func ChangeLog(close bool) {
	log.CloseLog = close
}

// GetOne get one data
func GetOne(columns []any, sql string, args ...any) (e error) {
	errmsg := ""
	start := time.Now()
	sql, e = interpolateParams(sql, args...)
	if e == nil {
		row := db.QueryRow(sql)
		e = row.Scan(columns...)
	}
	if e != nil {
		errmsg = e.Error()
	}
	cost := time.Since(start)
	if errmsg != "" {
		log.Error("%-12s\t[%s] err: %s", cost, sql, errmsg)
	} else {
		traceSlowSQL(sql, cost)
		log.Info("%-12s\t[%s]", cost, sql)
	}
	return
}

// GetOneStruct init to struct
func GetOneStruct(dest any, sql string, args ...any) (e error) {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to GetOneStruct dest")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to GetOneStruct dest")
	}
	if v.Elem().Type().Kind() != reflect.Struct {
		return errors.New("no struct passed to GetOneStruct dest")
	}
	rows, sql, e := list(sql, args...)

	isExists := false
	if e == nil {
		// start := time.Now()
		defer rows.Close()
		var columns []string
		columns, e = rows.Columns()
		if e == nil {
			for rows.Next() {
				e = rows.Scan(strutForScan(columns, dest)...)
				if e != nil {
					log.Error("[%s]数据赋值错误[%v]", sql, e)
					return
				}
				isExists = true
				break
			}
		}
		// cost := time.Since(start)
		// log.Info("%-12s\t解析[%s]", cost, sql)
	}
	if !isExists {
		e = errors.New("not exists")
	}

	return
}

// ListStruct init to struct list
func ListStruct(dest any, sql string, args ...any) (e error) {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to ListStruct dest")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to ListStruct dest")
	}
	if v.Elem().Type().Kind() != reflect.Slice {
		return errors.New("no slice passed to ListStruct dest")
	}

	rows, sql, e := list(sql, args...)
	if e == nil {
		// start := time.Now()
		defer rows.Close()
		sT := reflect.TypeOf(dest)
		// 取得数组中元素的类型
		sEE := sT.Elem().Elem()
		// 数组的值
		sVE := v.Elem()

		var columns []string
		columns, e = rows.Columns()
		if e == nil {
			for rows.Next() {
				vp := reflect.New(sEE)
				e = rows.Scan(strutForScan(columns, vp.Interface())...)
				if e != nil {
					log.Error("[%s]数据赋值错误[%v]", sql, e)
					return
				}
				resArr := reflect.Append(sVE, vp.Elem())
				sVE.Set(resArr)
			}

			if sVE.IsNil() {
				// 初始化数组，防止nil
				sVE.Set(reflect.MakeSlice(sT.Elem(), 0, 0))
			}
		}
		// cost := time.Since(start)
		// log.Info("%-12s\t解析[%s]", cost, sql)
	}
	return
}

// Insert insert value
func Insert(sql string, args ...any) (id int64, err error) {
	if db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		if err == nil {
			result, err = db.Exec(sql)
		}
		if err == nil {
			id, err = result.LastInsertId()
			if err != nil {
				errmsg = err.Error()
			}
		} else {
			errmsg = err.Error()
		}

		cost := time.Since(start)
		if errmsg != "" {
			log.Error("%-12s\t[%s] err: %s", cost, sql, errmsg)
		} else {
			traceSlowSQL(sql, cost)
			log.Info("%-12s\t[%s]", cost, sql)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Update update value
func Update(sql string, args ...any) (err error) {
	if db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		var num int64
		if err == nil {
			result, err = db.Exec(sql)
		}
		if err != nil {
			errmsg = err.Error()
		} else {
			num, err = result.RowsAffected()
			if err == nil && num < 0 {
				err = errors.New("update not working")
			}
		}
		cost := time.Since(start)
		if errmsg != "" {
			log.Error("%-12s\t[%s] rows affected = %d, err: %s", cost, sql, num, errmsg)
		} else {
			traceSlowSQL(sql, cost)
			log.Info("%-12s\t[%s] rows affected = %d", cost, sql, num)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Update2 更新数据
func Update2(sql string, args ...any) (rowsAffected int64, err error) {
	if db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		if err == nil {
			result, err = db.Exec(sql)
		}
		if err != nil {
			errmsg = err.Error()
		} else {
			rowsAffected, err = result.RowsAffected()
			if err == nil && rowsAffected < 0 {
				err = errors.New("update not working")
			}
		}
		cost := time.Since(start)
		if errmsg != "" {
			log.Error("%-12s\t[%s] rows affected = %d, err: %s", cost, sql, rowsAffected, errmsg)
		} else {
			traceSlowSQL(sql, cost)
			log.Info("%-12s\t[%s] rows affected = %d", cost, sql, rowsAffected)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Delete delete value
func Delete(sql string, args ...any) (rowsAffected int64, err error) {
	if db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		if err == nil {
			result, err = db.Exec(sql)
		}
		if err != nil {
			errmsg = err.Error()
		} else {
			rowsAffected, err = result.RowsAffected()
			if err == nil && rowsAffected < 0 {
				err = errors.New("delete not working")
			}
		}
		cost := time.Since(start)
		if errmsg != "" {
			log.Error("%-12s\t[%s] rows affected = %d, err: %s", cost, sql, rowsAffected, errmsg)
		} else {
			traceSlowSQL(sql, cost)
			log.Info("%-12s\t[%s] rows affected = %d,", cost, sql, rowsAffected)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Transaction 处理事务
func Transaction(fn func(*libSQL.Tx) bool) (e error) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		e = err
		return
	}

	if fn(tx) {
		e = tx.Commit()
	} else {
		e = tx.Rollback()
	}

	return
}
