package models

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/go-sql-driver/mysql"
	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/framework/database/mysqldb"
	"stove-gitlab.sginfra.net/backend/template/framework/database/redisdb"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
)

const (
	templateTable = "TB_TEMPLATE"

	errRedisNil           = "redis: nil"
	errNotFound           = "no rows in result set"
	errRedisPipeLineEmpty = "redis: pipeline is empty"
	errTemplateNotFound   = "template not found"

	sortRegAtDesc = "reg_at desc"
	sortRegAtAsc  = "reg_at asc"
)

func getTemplateTable() string {
	conf := config.GetInstance()
	return conf.DBPrefix + templateTable
}

// TemplateServiceDB Template 관련 구조체 정보
type TemplateServiceDB struct {
	conf    *config.TemplateConfig
	redisdb *redisdb.Redis
	mysqldb *mysqldb.Mysql
}

/*
	Public Functions
*/

// TemplateDB Template Mysql 통신을 위한 인터페이스
type TemplateDB interface {
	Init() error
	GetStatus() error

	CreateTable() error
	InitTable() error
	PingRedis() (bool, error)

	// Template Meta 관리
	CreateTemplate(template *Template) error
	UpdateTemplate(template *Template) error
	DeleteTemplate(template *Template) error
	GetTemplate(template *Template) error
	GetTemplates(templateList *TemplateList) error

	SetTemplateRedis(template *Template) error
	GetTemplateRedis(template *Template) error
	RemoveTemplateRedis(template *Template) error
	RefreshRedis(templateList *TemplateList) error
}

var templateDB TemplateDB
var onceDB sync.Once

func getDB() TemplateDB {
	onceDB.Do(func() {
		conf := config.GetInstance()
		templateDB = createDB(conf.DBType)
	})

	return templateDB
}

func createDB(dbtype string) TemplateDB {
	templateDB := &TemplateServiceDB{}
	if err := templateDB.Init(); err != nil {
		panic(err)
	}

	return templateDB
}

// Init 초기화 함수
func (o *TemplateServiceDB) Init() error {
	o.conf = config.GetInstance()

	var err error

	port := strconv.Itoa(o.conf.Database.RedisDB.Port)
	database := strconv.Itoa(o.conf.Database.RedisDB.DefaultDB)
	poolSize := strconv.Itoa(o.conf.Database.RedisDB.PoolSize)
	timeout := strconv.Itoa(o.conf.Database.RedisDB.TimeoutSec)
	if o.redisdb, err = redisdb.GetRedis(
		o.conf.Database.RedisDB.Host,
		port,
		o.conf.Database.RedisDB.Password,
		database,
		poolSize,
		timeout,
	); err != nil {
		log.Errorf("[Init] Redis err: %v", err)
		return err
	}

	if o.mysqldb, err = mysqldb.GetMysql(
		o.conf.Database.MysqlDB.Host,
		o.conf.Database.MysqlDB.ID,
		o.conf.Database.MysqlDB.Password,
		o.conf.Database.MysqlDB.Database,
		o.conf.Database.MysqlDB.PoolSize,
		o.conf.Database.MysqlDB.IdleSize,
	); err != nil {
		log.Errorf("[Init] Mysql err: %v", err)
		return err
	}

	return nil
}

// GetDBStatus DB 접속이 가능한지 확인한다.
func GetDBStatus() error {
	db := getDB()
	return db.GetStatus()
}

// GetStatus DB 접속 확인
func (o *TemplateServiceDB) GetStatus() error {
	connStr := mysql.NewConfig()
	connStr.Net = "tcp"
	connStr.Addr = o.conf.Database.MysqlDB.Host
	connStr.User = o.conf.Database.MysqlDB.ID
	connStr.Passwd = o.conf.Database.MysqlDB.Password
	connStr.DBName = o.conf.Database.MysqlDB.Database
	mysqlDsn := connStr.FormatDSN()

	db, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return err
	}
	db.Close()
	return nil
}

// CreateTable Template에 사용하는 Mysql의 Table들을 생성한다.
func CreateTable() error {
	db := getDB()
	return db.CreateTable()
}

// InitTable Template에 사용하는 모든 데이터를 삭제한다.
func InitTable() error {
	db := getDB()
	return db.InitTable()
}

// PingRedis Reids 통신 여부 체크
func PingRedis() (bool, error) {
	db := getDB()
	return db.PingRedis()
}

// CreateTable 테이블 생성
func (o *TemplateServiceDB) CreateTable() error {
	return o.CreateTemplateTable()
}

// InitTable 테이블 생성
func (o *TemplateServiceDB) InitTable() error {
	return o.InitTemplateTable()
}

// PingRedis Reids 통신 여부 체크
func (o *TemplateServiceDB) PingRedis() (bool, error) {
	return o.redisdb.Ping()
}
