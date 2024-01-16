package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	//"fmt"
	"time"
)

// 使用要应用 对应的 驱动引擎
// _ "github.com/mattn/go-sqlite3"
// _ "github.com/denisenkom/go-mssqldb"
// _"github.com/godror/godror"  oracledror?
// _ "github.com/go-sql-driver/mysql"
// _ "github.com/mattn/go-adodb"  //sql2000驱动
type SDatabase_KVList struct {
	link        *sql.DB
	maxtime     int
	maxconnum   int
	connections int
	isoracl     bool
}

// 初始化配置 连接前设置
// @ 最大时长 单位 为分钟
func (Class *SDatabase_KVList) ZPSetting(maxtime, maxconnum, connections int) {
	Class.maxtime = maxtime
	Class.maxconnum = maxconnum
	Class.connections = connections

}
func (Class *SDatabase_KVList) ZXSetHeartbeatTime(second int) {
	go DatabaseHeartbeat(Class.link, second)
}

func DatabaseHeartbeat(db *sql.DB, second int) {
	interval := second * 1000
	for {
		allTime.YCDelayProgram(interval)
		if err := db.Ping(); err != nil {
			fmt.Println("数据库断开:::", err.Error())
			return
		}
	}
}

func (Class *SDatabase_KVList) execSetting() {
	if Class.maxtime != 0 {
		Class.link.SetConnMaxLifetime(time.Minute * time.Duration(Class.maxtime))
	}
	if Class.maxconnum != 0 {
		Class.link.SetMaxOpenConns(Class.maxconnum)
	}
	if Class.connections != 0 {
		Class.link.SetMaxIdleConns(Class.connections)
	}

}

func (Class *SDatabase_KVList) CJCreateConectmssql(ip, port, dbname, username, password string) (returnerror error) {
	connectString := "server=" + ip + ";database=" + dbname + ";user id=" + username + ";password=" + password + ";port=" + port + ";encrypt=disable"
	Class.link, returnerror = sql.Open("mssql", connectString)
	if returnerror != nil {
		return
	}
	Class.execSetting()
	_, returnerror = Class.link.Exec("select name from SYSOBJECTS")
	return
}

func (Class *SDatabase_KVList) CJCreateLinkMssql2000(ip, port, dbname, username, password string) (returnerror error) {
	connectString := "Provider=SQLOLEDB;Initial Catalog=" + dbname + ";Data Source=" + ip + "," + port + ";user id=" + username + ";password=" + password
	Class.link, returnerror = sql.Open("adodb", connectString)
	if returnerror != nil {
		return
	}
	Class.execSetting()
	_, returnerror = Class.link.Exec("select name from SYSOBJECTS")
	return
}

func (Class *SDatabase_KVList) CJCreateLinkOracl(username, password, link string) (returnerror error) {
	connectString := "user=" + username + " password=" + password + " connectString=" + link
	Class.link, returnerror = sql.Open("godror", connectString)
	if returnerror != nil {
		return
	}
	Class.isoracl = true
	Class.execSetting()
	_, returnerror = Class.link.Exec("select * from user_tables")
	return
}

func (Class *SDatabase_KVList) CJCreateLinkSqlite3(ip string) (returnerror error) {
	Class.link, returnerror = sql.Open("sqlite3", ip)
	if returnerror != nil {
		return
	}
	Class.execSetting()
	_, returnerror = Class.link.Exec("SELECT name from sqlite_master where type='table'")
	return
}

func (Class *SDatabase_KVList) CJCreateLinkmysql(ip, port, dbname, username, password string) (returnerror error) {
	connectString := username + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbname
	Class.link, returnerror = sql.Open("mysql", connectString)
	if returnerror != nil {
		return
	}
	Class.execSetting()
	_, returnerror = Class.link.Query("SHOW TABLES")
	return
}
func (Class *SDatabase_KVList) CQuery(sql string, preventInjection ...bool) (list LList, returnerror error) {
	if len(preventInjection) > 0 && preventInjection[0] {
		isInjectText := strings.ToLower(sql)
		if strings.Contains(isInjectText, "delete") {
			returnerror = errors.New("不能包含关键字 delete")
			return
		}
		if strings.Contains(isInjectText, "insert") {
			returnerror = errors.New("不能包含关键字 insert")
			return
		}
		if strings.Contains(isInjectText, "update") {
			returnerror = errors.New("不能包含关键字 insert")
			return
		}

	}

	isoracl := Class.isoracl
	list.QClear()
	rows, err := Class.link.Query(sql)
	if err != nil {
		returnerror = err
		return
	}
	defer rows.Close()
	keyarr, _ := rows.Columns()
	cache := make([]any, len(keyarr))
	for index := range cache {
		var a any
		cache[index] = &a
	}
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]any)
		for i, data := range cache {
			v := *data.(*any)
			if isoracl {
				returnerror = KVListFilter(v)
				if returnerror != nil {
					v = any_to_doc(v)
				}
			}
			item[keyarr[i]] = v
			//item[键组[i]] = *data.(*any) //any中是什么 就返回什么
		}
		returnerror = list.TAddValue(item)
		if returnerror != nil {
			return
		}
	}

	return
}
func (Class *SDatabase_KVList) ZAdd(dbname string, table JKVtable) (returnerror error) {
	mapTable := table.Dtomap()
	keyarr := make([]string, 0)
	valuearr := make([]any, 0)
	replacearr := make([]string, 0)
	i := 0
	for k, v := range mapTable {
		switch nowvalue := v.(type) {
		case nil:
			continue
		case JKVtable:
			text := nowvalue.DtoJSON()
			keyarr = append(keyarr, text)
		case LList:
			text := nowvalue.DtoJSON()
			valuearr = append(valuearr, text)
		case []any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		case map[string]any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		default:
			valuearr = append(valuearr, nowvalue)
		}
		keyarr = append(keyarr, k)
		replacearr = append(replacearr, "?")
		i++
	}
	if len(keyarr) == 0 {
		returnerror = errors.New("错误:天加参数不能为空")
		return
	}
	sqlStr := "INSERT INTO " + dbname + " (" + strings.Join(keyarr, ",") + ")" + "VALUES (" + strings.Join(replacearr, ",") + ")"
	_, returnerror = Class.link.Exec(sqlStr, valuearr...)
	return
}

func (Class *SDatabase_KVList) GChange(dbname string, table JKVtable, condition ...string) (rows int, returnerror error) {
	mapTable := table.Dtomap()
	keyarr := make([]string, 0)
	valuearr := make([]any, 0)
	isoracl := Class.isoracl

	i := 0
	for k, v := range mapTable {
		switch nowvalue := v.(type) {
		case nil:
			continue
		case JKVtable:
			text := nowvalue.Dtomap()
			valuearr = append(valuearr, text)
		case LList:
			text := nowvalue.DtoJSON()
			valuearr = append(valuearr, text)
		case []any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		case map[string]any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		default:
			valuearr = append(valuearr, nowvalue)
		}
		if isoracl {
			keyarr = append(keyarr, k+"=:"+k)
		} else {
			keyarr = append(keyarr, k+"=?")
		}
		i++
	}
	sqlStr := ""
	if len(condition) >= 1 && condition[0] != "" {
		sqlStr = "UPDATE " + dbname + " SET " + strings.Join(keyarr, ",") + " WHERE " + condition[0]
	} else {
		sqlStr = "UPDATE " + dbname + " SET " + strings.Join(keyarr, ",")
	}

	ref, err := Class.link.Exec(sqlStr, valuearr...)
	if err != nil {
		returnerror = err
		return
	}
	n, _ := ref.RowsAffected()
	rows = int(n)
	return
}

func (Class *SDatabase_KVList) SDelete(dbname string, condition ...string) (rows int, returnerror error) {

	sqlStr := ""
	if len(condition) >= 1 && condition[0] != "" {
		sqlStr = "DELETE FROM  " + dbname + " WHERE " + condition[0]
	} else {
		sqlStr = "DELETE FROM  " + dbname
	}

	ref, err := Class.link.Exec(sqlStr)
	if err != nil {
		returnerror = err
		return
	}
	n, _ := ref.RowsAffected()
	rows = int(n)
	return
}
func (Class *SDatabase_KVList) ZXExecSql(sql string) (rows int, returnerror error) {
	ref, err := Class.link.Exec(sql)
	if err != nil {
		returnerror = err
		return
	}
	n, _ := ref.RowsAffected()
	rows = int(n)
	return
}

// 本对象 只能执行一次事务  S_事务提交 或 S_事务回滚 后 失效
func (Class *SDatabase_KVList) QSGetTransaction() (transaction SDatabase_KVList_Transaction, returnerror error) {
	conn, returnerror := Class.link.Begin()
	transaction.cInit(conn, Class.isoracl)
	return
}

type SDatabase_KVList_Transaction struct {
	conn    *sql.Tx
	isoracl bool
}

func (Class *SDatabase_KVList_Transaction) cInit(conn *sql.Tx, isoracl bool) {
	Class.conn = conn
	Class.isoracl = isoracl
}
func (Class *SDatabase_KVList_Transaction) SSubmit() (returnerror error) {
	returnerror = Class.conn.Commit()
	return
}
func (Class *SDatabase_KVList_Transaction) SRollback() (returnerror error) {
	returnerror = Class.conn.Rollback()
	return
}

func (Class *SDatabase_KVList_Transaction) ZAdd(dbname string, table JKVtable) (returnerror error) {
	mapTable := table.Dtomap()
	keyarr := make([]string, 0)
	valuearr := make([]any, 0)
	replacearr := make([]string, 0)
	i := 0
	for k, v := range mapTable {
		switch nowvalue := v.(type) {
		case nil:
			continue
		case JKVtable:
			text := nowvalue.DtoJSON()
			valuearr = append(valuearr, text)
		case LList:
			text := nowvalue.DtoJSON()
			valuearr = append(valuearr, text)
		case []any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		case map[string]any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		default:
			valuearr = append(valuearr, nowvalue)
		}
		keyarr = append(keyarr, k)
		replacearr = append(replacearr, "?")
		i++
	}
	if len(keyarr) == 0 {
		returnerror = errors.New("错误:天加参数不能为空")
		return
	}
	sqlStr := "INSERT INTO " + dbname + " (" + strings.Join(keyarr, ",") + ")" + "VALUES (" + strings.Join(replacearr, ",") + ")"
	_, returnerror = Class.conn.Exec(sqlStr, valuearr...)
	return
}

func (Class *SDatabase_KVList_Transaction) GChange(dbname string, table JKVtable, condition ...string) (rows int, returnerror error) {
	mapTable := table.Dtomap()
	keyarr := make([]string, 0)
	valuearr := make([]any, 0)
	isoracl := Class.isoracl
	i := 0
	for k, v := range mapTable {
		switch nowvalue := v.(type) {
		case nil:
			continue
		case JKVtable:
			text := nowvalue.DtoJSON()
			valuearr = append(valuearr, text)
		case LList:
			text := nowvalue.DtoJSON()
			valuearr = append(valuearr, text)
		case []any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		case map[string]any:
			valuearr = append(valuearr, any_to_doc(nowvalue))
		default:
			valuearr = append(valuearr, nowvalue)
		}
		if isoracl {
			keyarr = append(keyarr, k+"=:"+k)
		} else {
			keyarr = append(keyarr, k+"=?")
		}
		i++
	}
	sqlStr := ""
	if len(condition) >= 1 && condition[0] != "" {
		sqlStr = "UPDATE " + dbname + " SET " + strings.Join(keyarr, ",") + " WHERE " + condition[0]
	} else {
		sqlStr = "UPDATE " + dbname + " SET " + strings.Join(keyarr, ",")
	}

	ref, err := Class.conn.Exec(sqlStr, valuearr...)
	if err != nil {
		returnerror = err
		return
	}
	n, _ := ref.RowsAffected()
	rows = int(n)
	return
}

func (Class *SDatabase_KVList_Transaction) SDelete(dbname string, condition ...string) (rows int, returnerror error) {

	sqlStr := ""
	if len(condition) >= 1 && condition[0] != "" {
		sqlStr = "DELETE FROM  " + dbname + " WHERE " + condition[0]
	} else {
		sqlStr = "DELETE FROM  " + dbname
	}

	ref, err := Class.conn.Exec(sqlStr)
	if err != nil {
		returnerror = err
		return
	}
	n, _ := ref.RowsAffected()
	rows = int(n)
	return
}

func (Class *SDatabase_KVList_Transaction) ZXExecsql(sql string) (rows int, returnerror error) {
	ref, err := Class.conn.Exec(sql)
	if err != nil {
		returnerror = err
		return
	}
	n, _ := ref.RowsAffected()
	rows = int(n)
	return
}
