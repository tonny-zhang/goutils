package mysql

import (
	"context"
	libSQL "database/sql"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/tonny-zhang/goutils/logger"
)

// DB 数据库实例
type DB struct {
	Log logger.Logger
	db  *libSQL.DB
}

// 记录慢sql
func (ins DB) traceSlowSQL(sql string, cost time.Duration) {
	if cost > slowTime {
		fileNumInfo := ""
		_, file, line, _ := runtime.Caller(3)
		fileNumInfo = fmt.Sprintf("%s:%d", file, line)
		ins.Log.Warn("%s SLOW SQL %v>= %v [%s]", fileNumInfo, cost, slowTime, sql)
	}
}

// ChangeLog 更改日志
func (ins DB) ChangeLog(close bool) {
	ins.Log.CloseLog = close
}

// List 查询多条数据
func (ins DB) List(sql string, args ...any) (rows *libSQL.Rows, query string, err error) {
	if ins.db != nil {
		start := time.Now()

		sql, err = interpolateParams(sql, args...)
		if err == nil {
			rows, err = ins.db.Query(sql)
			// rows, err = db.Query(sql, args...)
			query = sql
		}
		cost := time.Since(start)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		if errmsg != "" {
			ins.Log.Error("%-12s\t[%s] err: %s", cost, sql, errmsg)
		} else {
			ins.traceSlowSQL(sql, cost)
			ins.Log.Info("%-12s\t[%s]", cost, sql)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}

	return
}

// GetOne get one data
func (ins DB) GetOne(columns []any, sql string, args ...any) (e error) {
	errmsg := ""
	start := time.Now()
	sql, e = interpolateParams(sql, args...)
	if e == nil {
		row := ins.db.QueryRow(sql)
		e = row.Scan(columns...)
	}
	if e != nil {
		errmsg = e.Error()
	}
	cost := time.Since(start)
	if errmsg != "" {
		ins.Log.Error("%-12s\t[%s] err: %s", cost, sql, errmsg)
	} else {
		ins.traceSlowSQL(sql, cost)
		ins.Log.Info("%-12s\t[%s]", cost, sql)
	}
	return
}

// GetOneStruct init to struct
func (ins DB) GetOneStruct(dest any, sql string, args ...any) (e error) {
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
	rows, sql, e := ins.List(sql, args...)

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
					ins.Log.Error("[%s]数据赋值错误[%v]", sql, e)
					return
				}
				isExists = true
				break
			}
		}
		// cost := time.Since(start)
		// ins.Log.Info("%-12s\t解析[%s]", cost, sql)
	}
	if !isExists {
		e = errors.New("not exists")
	}

	return
}

// ListStruct init to struct list
func (ins DB) ListStruct(dest any, sql string, args ...any) (e error) {
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

	rows, sql, e := ins.List(sql, args...)
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
					ins.Log.Error("[%s]数据赋值错误[%v]", sql, e)
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
		// ins.Log.Info("%-12s\t解析[%s]", cost, sql)
	}
	return
}

// Insert insert value
func (ins DB) Insert(sql string, args ...any) (id int64, err error) {
	if ins.db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		if err == nil {
			result, err = ins.db.Exec(sql)
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
			ins.Log.Error("%-12s\t[%s] err: %s", cost, sql, errmsg)
		} else {
			ins.traceSlowSQL(sql, cost)
			ins.Log.Info("%-12s\t[%s]", cost, sql)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Update update value
func (ins DB) Update(sql string, args ...any) (err error) {
	if ins.db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		var num int64
		if err == nil {
			result, err = ins.db.Exec(sql)
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
			ins.Log.Error("%-12s\t[%s] rows affected = %d, err: %s", cost, sql, num, errmsg)
		} else {
			ins.traceSlowSQL(sql, cost)
			ins.Log.Info("%-12s\t[%s] rows affected = %d", cost, sql, num)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Update2 更新数据
func (ins DB) Update2(sql string, args ...any) (rowsAffected int64, err error) {
	if ins.db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		if err == nil {
			result, err = ins.db.Exec(sql)
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
			ins.Log.Error("%-12s\t[%s] rows affected = %d, err: %s", cost, sql, rowsAffected, errmsg)
		} else {
			ins.traceSlowSQL(sql, cost)
			ins.Log.Info("%-12s\t[%s] rows affected = %d", cost, sql, rowsAffected)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Delete delete value
func (ins DB) Delete(sql string, args ...any) (rowsAffected int64, err error) {
	if ins.db != nil {
		start := time.Now()
		sql, err = interpolateParams(sql, args...)
		var result libSQL.Result

		errmsg := ""
		if err == nil {
			result, err = ins.db.Exec(sql)
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
			ins.Log.Error("%-12s\t[%s] rows affected = %d, err: %s", cost, sql, rowsAffected, errmsg)
		} else {
			ins.traceSlowSQL(sql, cost)
			ins.Log.Info("%-12s\t[%s] rows affected = %d,", cost, sql, rowsAffected)
		}
	} else {
		err = fmt.Errorf("db not inited")
	}
	return
}

// Transaction 处理事务
func (ins DB) Transaction(fn func(*libSQL.Tx) bool) (e error) {
	tx, err := ins.db.BeginTx(context.Background(), nil)
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
