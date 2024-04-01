package app

import (
	"flag"
	"fmt"
	"sync"

	"stove-gitlab.sginfra.net/backend/template/controllers/resultcode"

	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/controllers/externalapi"
	"stove-gitlab.sginfra.net/backend/template/controllers/internalapi"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	webconf "stove-gitlab.sginfra.net/backend/template/framework/web/config"
	"stove-gitlab.sginfra.net/backend/template/models"
)

var (
	command     = flag.String("cmd", "", "cmd")
	serviceType = flag.String("service_type", "", "service_type")
	serviceID   = flag.String("service_id", "", "service_id")
)

// TemplateApp Template 서버 구동을 위한 구조체
type TemplateApp struct {
	base.BaseApp
	conf       *config.TemplateConfig
	configFile string
}

// Init TemplateApp Init
func (o *TemplateApp) Init(configFile string) error {
	// 싱글턴 객체 초기화
	o.conf = config.GetInstance(configFile)
	//o.conf.Version = o.Version

	// 결과코드 등록
	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	// API Request Paramter 등록
	context.AppendRequestParameter()

	return nil
}

// Run TemplateApp Run
func (o *TemplateApp) Run(wg *sync.WaitGroup) error {
	return nil
}

// GetConfig TemplateApp Config 로드
func (o *TemplateApp) GetConfig() *webconf.Config {
	return &o.conf.Config
}

// CreateTable ./Template -cli -cmd=createtable -c config.yml
func (o *TemplateApp) CreateTable() error {
	if err := models.CreateTable(); err != nil {
		return err
	}

	return nil
}

// InitTable ./Template -cli -cmd=inittable -c config.yml
func (o *TemplateApp) InitTable() error {
	if err := models.InitTable(); err != nil {
		return err
	}

	return nil
}

// Execute TemplateApp Execute
func (o *TemplateApp) Execute() error {

	if *command == "createtable" {
		fmt.Println("createtable start")
		o.CreateTable()
		fmt.Println("createtable end")
		return nil
	}

	if *command == "inittable" {
		fmt.Println("inittable start")
		o.InitTable()
		fmt.Println("inittable end")
		return nil
	}

	return nil
}

// NewApp TemplateApp NewApp
func NewApp(version string) (*TemplateApp, error) {
	app := &TemplateApp{}
	//app.Version = version
	intAPI := internalapi.NewAPI()
	extAPI := externalapi.NewAPI()

	// 실행 매개변수를 추가하려면 이 곳에서 구현
	// test := flag.String("t", "hello", "test parameter")
	// fmt.Println(*test)

	if err := app.BaseApp.Init(app, intAPI, extAPI); err != nil {
		return nil, err
	}

	return app, nil
}
