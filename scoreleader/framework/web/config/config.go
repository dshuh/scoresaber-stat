package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
	"stove-gitlab.sginfra.net/backend/template/framework/database/mongodb"
	"stove-gitlab.sginfra.net/backend/template/framework/database/mysqldb"
	"stove-gitlab.sginfra.net/backend/template/framework/database/redisdb"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
)

const (
	// ErrInitConfigFailed Config를 위한 Singleton 인스턴스 생성 실패 에러
	ErrInitConfigFailed = "Can't make a singleton instance for app_config"
)

// Config 는 어플리케이션의 설정을 처리한다.
type Config struct {
	// Auth        Auth    `json:"auth" yaml:"auth"`
	// Aws         Aws     `json:"aws" yaml:"aws"`
	Log         log.Log  `json:"log" yaml:"log"`
	Environment string   `json:"environment" yaml:"environment"` // dev, qa, sandbox, live
	DBType      string   `json:"db_type" yaml:"db_type"`         // dynamodb, mysql, couchbase
	DBPrefix    string   `json:"db_prefix" yaml:"db_prefix"`
	Database    Database `json:"database" yaml:"database"`
	// Cache       redisdb.RedisConfig `json:"cache_option" yaml:"cache_option"`
	APIServers []APIServer `json:"api_servers" yaml:"api_servers"`
}

// APIServer API 서버 설정
type APIServer struct {
	Name              string    `json:"name" yaml:"name"`
	Port              int       `json:"port" yaml:"port"`
	Disable           bool      `json:"disable" yaml:"disable"`
	Routes            string    `json:"routes" yaml:"routes"`
	ReadTimeOut       int       `json:"read_timeout" yaml:"read_timeout"`
	ReadHeaderTimeout int       `json:"read_header_timeout" yaml:"read_header_timeout"`
	WriteTimeout      int       `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout       int       `json:"idle_timeout" yaml:"idle_timeout"`
	DetailLog         DetailLog `json:"detail_log" yaml:"detail_log"`
	DisableCORS       bool      `json:"disable_cors" yaml:"disable_cors"`
}

// DetailLog 상세 로그 기록 설정 (RequestBody, ResponseBody 기록 설정)
type DetailLog struct {
	Enable              bool       `json:"enable" yaml:"enable"`
	Request             ReqRespLog `json:"request" yaml:"request"`
	Response            ReqRespLog `json:"response" yaml:"response"`
	LatencyUnitMilliSec bool       `json:"latency_unit_millisec" yaml:"latency_unit_millisec"`
}

// ReqRespLog Request, Response 상세 로그 설정
type ReqRespLog struct {
	Enable bool `json:"enable" yaml:"enable"`
	MaxLen int  `json:"max_len" yaml:"max_len"`
}

// Database DB 설정 정보
type Database struct {
	MysqlDB mysqldb.MysqlConfig `json:"mysql" yaml:"mysql"`
	MongoDB mongodb.MongoConfig `json:"mongo" yaml:"mongo"`
	RedisDB redisdb.RedisConfig `json:"redis" yaml:"redis"`
}

type AuthDepAPI struct {
	GetMemberByAccessToken  string `json:"get_member_by_access_token" yaml:"get_member_by_access_token"`
	GetAccessToken          string `json:"get_access_token" yaml:"get_access_token"`
	EncryptPwd              string `json:"encrypt_pwd" yaml:"encrypt_pwd"`
	GetApplication          string `json:"get_application" yaml:"get_application"`
	GetApplicationByNo      string `json:"get_application_by_no" yaml:"get_application_by_no"`
	PostApplication         string `json:"post_application" yaml:"post_application"`
	GetApplicationCallback  string `json:"get_application_callback" yaml:"get_application_callback"`
	PostApplicationCallback string `json:"post_application_callback" yaml:"post_application_callback"`
}

// // Auth 는 플랫폼 인증을 위한 정보를 설정한다.
// type Auth struct {
// 	JWTSecret string     `json:"jwt_secret" yaml:"jwt_secret"`
// 	DepAPI    AuthDepAPI `json:"dep_api" yaml:"dep_api"`
// }

// // Aws AWS 관련 설정
// type Aws struct {
// 	Region      string `json:"region" yaml:"region"`
// 	Credentials string `json:"credentials" yaml:"credentials"`
// }

// Load 함수는 설정파일에서 설정을 불러온다.
func Load(filepath string, o interface{}) error {
	ext := path.Ext(filepath)

	if ext == ".json" {
		return loadJSONConfig(filepath, o)
	}

	if ext == ".yml" {
		return loadYMLConfig(filepath, o)
	}

	return errors.New("This file format is not supported")
}

func loadJSONConfig(filepath string, o interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&o)
	if err != nil {
		return err
	}

	return nil
}

func loadYMLConfig(filepath string, o interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, o)
}
