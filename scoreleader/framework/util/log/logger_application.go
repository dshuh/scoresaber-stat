package log

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"

	"stove-gitlab.sginfra.net/backend/template/framework/util/webhook"
)

// Log 로그 컨피그 설정 구조체
//
// ServiceName: 서비스명, 로그에 servicename 필드에 출력되는 정보
// Path: 로그 파일 경로
// Rotation: 파일 로테이션 기능 활성화
// Level: 로그 출력 레벨 설정("panic", "fatal", "error", "warning", "info", "debug", "trace")
// Formatter: 로그 포맷 설정("json", "text")
// Console: 로그 콘솔 출력 활성화
// WithModule: 로그 출력 시 모듈 정보 추가 활성화
// WithFunc: 로그 출력 시 함수 정보 추가 활성화
// WithLine: 로그 출력 시 라인 정보 추가 활성화
// APITrace: 안정화 지표 로그 설정 구조체
// CustomField: 로그 필드 이름 변경 설정 구조체
type Log struct {
	ServiceName  string                   `json:"service_name" yaml:"service_name"`
	Path         string                   `json:"path" yaml:"path"`
	Rotation     bool                     `json:"rotation" yaml:"rotation"`
	Level        string                   `json:"level" yaml:"level"`
	Formatter    string                   `json:"formatter" yaml:"formatter"`
	Console      bool                     `json:"console" yaml:"console"`
	WithModule   bool                     `json:"with_module" yaml:"with_module"`
	WithFunc     bool                     `json:"with_func" yaml:"with_func"`
	WithLine     bool                     `json:"with_line" yaml:"with_line"`
	APITrace     APITrace                 `json:"apitrace" yaml:"apitrace"`
	CustomField  CustomField              `json:"custom_field" yaml:"custom_field"`
	SmilehubHook webhook.SmilehubHookConf `json:"smilehub_hook" yaml:"smilehub_hook"`
}

// CustomField 는 출력되는 로그 필드의 이름을 원하는대로 변경할 수 있도록 한다.
type CustomField struct {
	Use   bool   `json:"use" yaml:"use"`
	Time  string `json:"time" yaml:"time"`
	Level string `json:"level" yaml:"level"`
	Msg   string `json:"msg" yaml:"msg"`
	Func  string `json:"func" yaml:"func"`
}

// Logger logger 구조체
type Logger struct {
	appLogger *logrus.Logger
	conf      *Log
	apiLogger *logrus.Logger
}

var gLogger *Logger
var gOnceLogger sync.Once

// GetLogger logger 의 싱글턴, 로거 초기화 및 반환 함수
func GetLogger(params ...interface{}) (*logrus.Logger, error) {
	gOnceLogger.Do(func() {
		if len(params) <= 0 {
			panic("Logger config parameter is invalid")
		}
		conf, ok := params[0].(*Log)
		if !ok {
			panic("Logger config parameter is invalid")
		}
		logrus.SetFormatter(&logrus.JSONFormatter{})

		gLogger = &Logger{
			appLogger: logrus.StandardLogger(),
			conf:      conf,
		}

		// 플랫폼운영팀 요청 사항. 배포시 보안이슈로 인해 644 권한 변경
		// Log Path가 비어 있다면, 파일출력은 하지 않음
		if gLogger.conf.Path == "" {
			logrus.SetOutput(os.Stdout)
		} else {
			// create a log directory
			dir, _ := filepath.Split(gLogger.conf.Path)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					logrus.Error(err)
				}
			}

			if !gLogger.conf.Rotation {
				f, err := os.OpenFile(gLogger.conf.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					logrus.Error("Log path is [", gLogger.conf.Path, "],", err)
				}

				if gLogger.conf.Console {
					logrus.SetOutput(io.MultiWriter(f, os.Stdout))
				} else {
					logrus.SetOutput(f)
				}
			} else {
				// 플랫폼운영팀 요청 사항. 로그 로테이션 추가
				// NOTE: rotation 기능에 대한 안정성 추가 필요
				logpath := gLogger.conf.Path
				if !filepath.IsAbs(logpath) {
					pwd, _ := os.Getwd()
					logpath = pwd + string(os.PathSeparator) + gLogger.conf.Path
				}

				// rotation time : 86400 sec
				// max age : 7 day
				writer, _ := rotatelogs.New(
					logpath+".%Y%m%d",
					rotatelogs.WithLinkName(logpath+".log"),
					rotatelogs.WithRotationTime(time.Hour*24),
					rotatelogs.WithClock(rotatelogs.UTC),
				)

				if gLogger.conf.Console {
					logrus.SetOutput(io.MultiWriter(writer, os.Stdout))
				} else {
					logrus.SetOutput(writer)
				}
			}

			level, err := logrus.ParseLevel(gLogger.conf.Level)
			if err != nil {
				level = logrus.InfoLevel
			}
			logrus.SetLevel(level)

			if gLogger.conf.Formatter == "json" {
				jsonFormatter := &logrus.JSONFormatter{}
				jsonFormatter.TimestampFormat = time.RFC3339Nano
				if gLogger.conf.CustomField.Use {
					jsonFormatter.FieldMap = getCustomField()
				}
				logrus.SetFormatter(jsonFormatter)
			} else {
				textFormatter := &logrus.TextFormatter{}
				textFormatter.TimestampFormat = time.RFC3339Nano
				if gLogger.conf.CustomField.Use {
					textFormatter.FieldMap = getCustomField()
				}
				logrus.SetFormatter(textFormatter)
			}
		}

		if gLogger.conf.APITrace.Use {
			// new APITraceLogger
			if err := gLogger.initAPITraceLogger(); err != nil {
				panic(err)
			}
		}

		// smile hub webhook
		if conf.SmilehubHook.Use {
			u, err := url.Parse(conf.SmilehubHook.Url)
			if err != nil {
				panic(err)
			}
			hook, err := webhook.NewSmilehubHook(u, []logrus.Level{logrus.ErrorLevel}, &logrus.TextFormatter{}, conf.SmilehubHook.WorkerSize)
			if err != nil {
				panic(err)
			}
			gLogger.appLogger.Hooks.Add(hook)
		}
	})

	if len(params) == 0 {
		if gLogger == nil {
			return nil, errors.New("logger is not initialized")
		}
		return gLogger.appLogger, nil
	}

	return gLogger.appLogger, nil
}

// initAPITraceLogger 안정화 지표 용 로거 초기화
func (l *Logger) initAPITraceLogger() error {
	// new APITraceLogger
	apiTraceConf := gLogger.conf.APITrace
	logpath := apiTraceConf.Path

	// create a log directory
	dir, _ := filepath.Split(logpath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if !filepath.IsAbs(logpath) {
		pwd, _ := os.Getwd()
		logpath = pwd + string(os.PathSeparator) + logpath
	}

	var rotateTime = 24 * time.Hour
	var maxAge = 7 * 24 * time.Hour

	if apiTraceConf.RotateTime > 0 {
		rotateTime = time.Duration(apiTraceConf.RotateTime) * time.Hour
	}

	if apiTraceConf.KeepDay > 0 {
		maxAge = time.Duration(apiTraceConf.KeepDay) * 24 * time.Hour
	}

	writer, err := rotatelogs.New(
		logpath+".%Y%m%d",
		rotatelogs.WithLinkName(logpath+".log"),
		rotatelogs.WithRotationTime(rotateTime),
		rotatelogs.WithClock(rotatelogs.UTC),
		rotatelogs.WithMaxAge(maxAge),
	)
	if err != nil {
		panic(err)
	}

	l.apiLogger = logrus.New()
	l.apiLogger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})

	if apiTraceConf.Console {
		l.apiLogger.SetOutput(io.MultiWriter(writer, os.Stdout))
	} else {
		l.apiLogger.SetOutput(writer)
	}

	return nil
}

// GetAPITraceLogger 안정화 지표를 위한 파일 로거 반환
func GetAPITraceLogger() *logrus.Logger {
	return gLogger.apiLogger
}

func getLogFields() *logrus.Fields {
	return getLogFieldsWithSkip(4)
}

func getLogFieldsWithSkip(callerSkip int) *logrus.Fields {
	logFields := logrus.Fields{}
	if gLogger.conf.ServiceName != "" {
		logFields["servicename"] = gLogger.conf.ServiceName
	}
	if pc, file, line, ok := runtime.Caller(callerSkip); ok {
		funname := runtime.FuncForPC(pc).Name()

		if gLogger.conf.WithModule {
			logFields["file"] = file
		}
		if gLogger.conf.WithFunc {
			logFields["func"] = funname
		}
		if gLogger.conf.WithLine {
			logFields["line"] = line
		}
		logFields["trace"] = getTrace(callerSkip + 2)
	}
	return &logFields
}

// func getAPITraceLoggerWithRuntimeContext() *logrus.Entry {
// 	if logFields := getLogFields(); logFields != nil {
// 		return gLogger.apiLogger.WithFields(*logFields)
// 	} else {
// 		return nil
// 	}
// }

func getLoggerWithRuntimeContextWithSkip(skip int) *logrus.Entry {
	if logFields := getLogFieldsWithSkip(skip); logFields != nil {
		return gLogger.appLogger.WithFields(*logFields)
	} else {
		return nil
	}
}

func getLoggerWithRuntimeContext() *logrus.Entry {
	if logFields := getLogFields(); logFields != nil {
		return gLogger.appLogger.WithFields(*logFields)
	} else {
		return nil
	}
}

// Debug log api
func Debug(args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Debug(args...)
	} else {
		gLogger.appLogger.Debug(args...)
	}
}

// Info log api
func Info(args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Info(args...)
	} else {
		gLogger.appLogger.Info(args...)
	}
}

// Warn log api
func Warn(args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Warn(args...)
	} else {
		gLogger.appLogger.Warn(args...)
	}
}

// Error log api
func Error(args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Error(args...)
	} else {
		gLogger.appLogger.Error(args...)
	}
}

func ErrorWithSkip(skip int, args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContextWithSkip(skip); entry != nil {
		entry.Error(args...)
	} else {
		gLogger.appLogger.Error(args...)
	}
}

// Fatal log api
func Fatal(args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Fatal(args...)
	} else {
		gLogger.appLogger.Fatal(args...)
	}
}

// Panic log api
func Panic(args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Panic(args...)
	} else {
		gLogger.appLogger.Panic(args...)
	}
}

// WithError 에러 객체를 포함하여 로그 출력
func WithError(err error, args ...interface{}) {
	if gLogger == nil {
		fmt.Println("error: ", err, args)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.WithError(err).Error(args...)
	} else {
		gLogger.appLogger.WithError(err).Error(args...)
	}
}

// Debugf log api
func Debugf(format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Debugf(format, args...)
	} else {
		gLogger.appLogger.Debugf(format, args...)
	}
}

// Infof log api
func Infof(format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Infof(format, args...)
	} else {
		gLogger.appLogger.Infof(format, args...)
	}
}

// Warnf log api
func Warnf(format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Warnf(format, args...)
	} else {
		gLogger.appLogger.Warnf(format, args...)
	}
}

// Errorf log api
func Errorf(format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Errorf(format, args...)
	} else {
		gLogger.appLogger.Errorf(format, args...)
	}
}

// Fatalf log api
func Fatalf(format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Fatalf(format, args...)
	} else {
		gLogger.appLogger.Fatalf(format, args...)
	}
}

// Panicf log api
func Panicf(format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf(format+"\n", args...)
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.Panicf(format, args...)
	} else {
		gLogger.appLogger.Panicf(format, args...)
	}
}

// WithErrorf 에러 객체를 포함하여 로그 출력
func WithErrorf(err error, format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Printf("error: %v, "+format+"\n", err, args)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.WithError(err).Errorf(format, args...)
	} else {
		gLogger.appLogger.WithError(err).Errorf(format, args...)
	}
}

func getTrace(callerSkip int) string {
	const traceCnt int = 5
	var trace string
	for end := callerSkip + traceCnt; callerSkip < end; callerSkip++ {
		// 루프 돌면서 콜 스택 저장
		if pc, _, line, ok := runtime.Caller(callerSkip); ok {
			funname := runtime.FuncForPC(pc).Name()
			trace += fmt.Sprintf("%v:%v ", funname, line)
		} else {
			break
		}
	}
	return trace
}

func getCustomField() logrus.FieldMap {
	fieldMap := logrus.FieldMap{}

	if len(gLogger.conf.CustomField.Time) > 0 {
		fieldMap[logrus.FieldKeyTime] = gLogger.conf.CustomField.Time
	}
	if len(gLogger.conf.CustomField.Level) > 0 {
		fieldMap[logrus.FieldKeyLevel] = gLogger.conf.CustomField.Level
	}
	if len(gLogger.conf.CustomField.Msg) > 0 {
		fieldMap[logrus.FieldKeyMsg] = gLogger.conf.CustomField.Msg
	}
	if len(gLogger.conf.CustomField.Func) > 0 {
		fieldMap[logrus.FieldKeyFunc] = gLogger.conf.CustomField.Func
	}
	return fieldMap
}

type Fields map[string]interface{}

func ErrorWithField(fields Fields, args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.WithFields(logrus.Fields(fields)).Error(args...)
	} else {
		gLogger.appLogger.WithFields(logrus.Fields(fields)).Error(args...)
	}
}

func ErrorfWithField(fields Fields, format string, args ...interface{}) {
	if gLogger == nil {
		fmt.Println(args...)
		return
	}
	if entry := getLoggerWithRuntimeContext(); entry != nil {
		entry.WithFields(logrus.Fields(fields)).Errorf(format, args...)
	} else {
		gLogger.appLogger.WithFields(logrus.Fields(fields)).Errorf(format, args...)
	}
}
