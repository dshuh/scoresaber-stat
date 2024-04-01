package base

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
	"stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

var AppVersion string

type ICleanUp interface {
	CleanUp()
}

// ICLI Command Line Interface
type ICLI interface {
	Execute() error
}

type IApp interface {
	Init(configFile string) error
	Run(wg *sync.WaitGroup) error
	GetConfig() *config.Config
}

type BaseApp struct {
	lock sync.Mutex

	APIControllers []IController
	App            IApp

	baseConf   *config.Config
	configFile string
	cli        bool // 명령줄 인터페이스로 사용할지 여부
	version    bool

	isRunning bool
	wg        sync.WaitGroup

	sigch  chan os.Signal
	flowch []chan bool
}

func (o *BaseApp) IsRunning() bool {
	return o.isRunning
}

func (o *BaseApp) Init(app IApp, apiControllers ...IController) error {
	o.App = app
	o.APIControllers = apiControllers
	o.isRunning = false

	var configFile string
	var cpuNum int
	var version, cli bool
	if flag.Lookup("c") == nil {
		flag.StringVar(&configFile, "c", "config.yml", "Config file path")
	}
	configFile = flag.Lookup("c").Value.(flag.Getter).Get().(string)

	if flag.Lookup("cpu") == nil {
		flag.IntVar(&cpuNum, "cpu", runtime.GOMAXPROCS(-1), "The number of cpus")
	}
	cpuNum = flag.Lookup("cpu").Value.(flag.Getter).Get().(int)

	if flag.Lookup("v") == nil {
		flag.BoolVar(&version, "v", false, "Print version")
	}
	version = flag.Lookup("v").Value.(flag.Getter).Get().(bool)

	if flag.Lookup("cli") == nil {
		flag.BoolVar(&cli, "cli", false, "Use as CLI")
	}
	cli = flag.Lookup("cli").Value.(flag.Getter).Get().(bool)

	flag.Parse()

	runtime.GOMAXPROCS(cpuNum)

	o.configFile = configFile
	o.version = version
	o.cli = cli

	return nil
}

func (o *BaseApp) SetCLI(enable bool) {
	o.cli = enable
}

func (o *BaseApp) Interrupt() {
	o.sigch <- os.Interrupt
}

func (o *BaseApp) Stop() {
	o.wg.Wait()
	if cu, ok := o.App.(ICleanUp); ok {
		cu.CleanUp()
	}

	log.Info("Application is terminated.")
}

func (o *BaseApp) Start() error {
	// 어플리케이션을 관리명령어 모드로 실행
	// - 버전 출력
	if o.version {
		fmt.Println(os.Args[0], AppVersion)
		return nil
	}

	if err := o.App.Init(o.configFile); err != nil {
		return err
	}

	o.baseConf = o.App.GetConfig()
	if o.baseConf == nil {
		return errors.New("Base config is not initialized")
	}

	// 싱글턴 객체 초기화
	_, err := log.GetLogger(&o.baseConf.Log)
	if err != nil {
		return err
	}

	if o.baseConf.DBType != "" {
		if o.baseConf.DBType == "mongodb" {
			log.Debug("DBType:", o.baseConf.DBType)
			//if _, err := amigodb.GetMongoDB(o.baseConf.DBAuth.Host, o.baseConf.DBAuth.ID, o.baseConf.DBAuth.Password, o.baseConf.DBAuth.Database); err != nil {
			//	log.Error("MongoDB connection failed")
			//	return err
			//}
		} else if o.baseConf.DBType == "redis" {
			// if sess, err := amigodb.GetRedis(o.baseConf.DBAuth.Host, o.baseConf.DBAuth.Port, o.baseConf.DBAuth.Password, o.baseConf.DBAuth.Database, o.baseConf.DBAuth.PoolSize); err != nil || sess == nil {
			// 	log.Error("Redis connection failed")
			// 	return err
			// }
		} else if o.baseConf.DBType == "mysql" {
			// if _, err := amigodb.GetMysql(o.baseConf.DBAuth.Host, o.baseConf.DBAuth.ID, o.baseConf.DBAuth.Password, o.baseConf.DBAuth.Database); err != nil {
			// 	log.Error("Mysql connection failed")
			// 	return err
			// }
		} else if o.baseConf.DBType == "dynamodb" {
			// if _, err := amigoaws.GetDynamoDB(o.baseConf.Aws.Region, o.baseConf.Aws.Credentials); err != nil {
			// 	log.Error("DynamoDB connection failed")
			// 	return err
			// }
		} else {
			return errors.New("[" + o.baseConf.DBType + "] database type is not supported")
		}
	}

	// - 콘솔 커맨드
	if cc, ok := o.App.(ICLI); ok && o.cli {
		if err := cc.Execute(); err != nil {
			return err
		}
		return nil
	}

	o.flowch = make([]chan bool, 0)
	o.sigch = make(chan os.Signal, 1)
	signal.Notify(o.sigch, syscall.SIGINT, syscall.SIGTERM)

	// 어플리케이션 시작
	for _, apiController := range o.APIControllers {
		if apiController == nil {
			continue
		}
		if apiController.GetConfig().Disable {
			continue
		}
		o.wg.Add(1)
		go o.runAPIServer(apiController)
	}

	if err := o.App.Run(&o.wg); err != nil {
		return err
	}

	// OS 인터럽트 시그널 대기
	o.isRunning = true
	<-o.sigch
	for _, c := range o.flowch {
		c <- true
	}

	o.isRunning = false

	o.Stop()
	return nil
}

func (o *BaseApp) runAPIServer(controller IController) {
	defer o.wg.Done()

	conf := controller.GetConfig()
	logger, err := log.GetLogger()
	if err != nil {
		log.Error(err.Error())
		return
	}
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	if conf.DetailLog.Enable {
		e.Use(DetailInfoLogger(o.baseConf.Log.ServiceName, conf.Name, &conf.DetailLog, logger))
	} else {
		e.Use(NewWithNameAndLogger(o.baseConf.Log.ServiceName, conf.Name, logger))
	}
	if o.baseConf.Log.APITrace.Use {
		e.Use(APITraceLogger(o.baseConf.Log.APITrace.ExcludeUrl))
	}

	//CORS
	if !conf.DisableCORS {
		e.Use(middleware.CORS())
	}

	if err := controller.Init(e); err != nil {
		log.Error(err.Error())
		return
	}

	go func() {
		echoStopch := make(chan bool)
		o.lock.Lock()
		o.flowch = append(o.flowch, echoStopch)
		o.lock.Unlock()
		<-echoStopch
		e.Close()
	}()

	// Server
	addr := fmt.Sprintf(":%d", conf.Port)

	readTimeout := 90
	if conf.ReadTimeOut > 0 {
		readTimeout = conf.ReadTimeOut
	}
	e.Server.ReadTimeout = time.Duration(readTimeout) * time.Second
	log.Infof("%s started, read timeout set %dsec", conf.Name, readTimeout)

	idleTimeout := 90
	if conf.IdleTimeout > 0 {
		idleTimeout = conf.IdleTimeout
	}
	e.Server.IdleTimeout = time.Duration(idleTimeout) * time.Second
	log.Infof("%s started, idle timeout set %dsec", conf.Name, idleTimeout)

	if conf.ReadHeaderTimeout > 0 {
		e.Server.ReadHeaderTimeout = time.Duration(conf.ReadHeaderTimeout) * time.Second
		log.Infof("%s started, read header timeout set %dsec", conf.Name, conf.ReadHeaderTimeout)
	}
	if conf.WriteTimeout > 0 {
		e.Server.WriteTimeout = time.Duration(conf.WriteTimeout) * time.Second
		log.Infof("%s started, write timeout set %dsec", conf.Name, conf.WriteTimeout)
	}

	if err := e.Start(addr); err != nil {
		log.Error(err)
	}
}
