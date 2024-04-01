package mysqldb

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"sync"

	"github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db *sql.DB
}

// MysqlDB Mysql 설정 정보
type MysqlConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	ID       string `json:"id" yaml:"id"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	PoolSize string `json:"poolsize" yaml:"poolsize"`
	IdleSize string `json:"idlesize" yaml:"idlesize"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	ReadPref string `json:"readpref" yaml:"readpref"`
}

var gMysql *Mysql
var gOnceMysql sync.Once

func GetMysql(params ...string) (*Mysql, error) {
	var err error
	err = nil

	gOnceMysql.Do(func() {
		host := ""
		user := ""
		passwd := ""
		dbname := ""
		poolSize := 0
		idleSize := 0

		if len(params) >= 1 {
			host = params[0]
		}
		if len(params) >= 2 {
			user = params[1]
		}
		if len(params) >= 3 {
			passwd = params[2]
		}
		if len(params) >= 4 {
			dbname = params[3]
		}
		if len(params) >= 5 {
			poolSize, err = strconv.Atoi(params[4])
			if err != nil {
				poolSize = 0
			}
		}
		if len(params) >= 6 {
			idleSize, err = strconv.Atoi(params[5])
			if err != nil {
				idleSize = 0
			}
		}

		connStr := mysql.NewConfig()
		connStr.User = user
		connStr.Passwd = passwd
		connStr.Net = "tcp"
		connStr.Addr = host
		connStr.DBName = dbname
		mysqlDsn := connStr.FormatDSN()

		gMysql, err = NewMysql(mysqlDsn, poolSize, idleSize)
	})

	if err != nil {
		gMysql = nil
		return nil, err
	}

	if len(params) == 0 {
		if gMysql == nil {
			return nil, errors.New("DB is not initialized")
		}
		return gMysql, nil
	}

	return gMysql, nil
}

func NewMysql(connStr string, poolSize int, idleSize int) (mydb *Mysql, err error) {
	mydb = new(Mysql)
	mydb.db, err = sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	err = mydb.db.Ping()
	if err != nil {
		return nil, err
	}

	if poolSize > 0 {
		mydb.db.SetMaxOpenConns(poolSize)
	}
	if idleSize > 0 {
		mydb.db.SetMaxIdleConns(idleSize)
	}

	return mydb, err
}

func (o *Mysql) Excute(sql string) (sql.Result, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}
	return o.db.Exec(sql)
}

func (o *Mysql) ExcuteWithArgs(query string, args ...interface{}) (sql.Result, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}
	return o.db.Exec(query, args...)
}

func (o *Mysql) QueryRow(sql string, column ...interface{}) error {
	if o.db == nil {
		return errors.New("DB is nil")
	}

	return o.db.QueryRow(sql).Scan(column...)
}

func (o *Mysql) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	if args != nil {
		return o.db.Query(sql, args...)
	} else {
		return o.db.Query(sql)
	}
}

func (o *Mysql) QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	if args != nil {
		return o.db.QueryContext(ctx, sql, args...)
	} else {
		return o.db.QueryContext(ctx, sql)
	}
}

func (o *Mysql) PrepareAndQuery(query string, args ...interface{}) (*sql.Rows, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	stmt, err := o.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if args != nil {
		return stmt.Query(args...)
	} else {
		return stmt.Query()
	}
}

func (o *Mysql) PrepareAndExec(query string, args ...interface{}) (sql.Result, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	stmt, err := o.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// transaction 관련 api
func (o *Mysql) Begin() (*sql.Tx, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	return o.db.Begin()
}

func (o *Mysql) TxPrepareAndExec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *Mysql) TxQuery(tx *sql.Tx, sql string, args ...interface{}) (*sql.Rows, error) {
	if o.db == nil {
		return nil, errors.New("DB is nil")
	}

	if args != nil {
		return tx.Query(sql, args...)
	} else {
		return tx.Query(sql)
	}

}

func (o *Mysql) TxQueryRow(tx *sql.Tx, sql string, column ...interface{}) error {
	if o.db == nil {
		return errors.New("DB is nil")
	}

	return tx.QueryRow(sql).Scan(column...)
}

func (o *Mysql) Stats() (sql.DBStats, error) {
	if o.db == nil {
		return sql.DBStats{}, errors.New("DB is nil")
	}
	return o.db.Stats(), nil
}
