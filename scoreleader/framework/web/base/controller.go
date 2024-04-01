package base

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
	"stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

type IController interface {
	Init(e *echo.Echo) error
	GetConfig() *config.APIServer
}

type Route struct {
	Handler        string   `json:"handler" yaml:"handler"`
	Method         string   `json:"method" yaml:"method"`
	URI            string   `json:"uri" yaml:"uri"`
	MinVersion     string   `json:"minversion" yaml:"minversion"`
	MaxVersion     string   `json:"maxversion" yaml:"maxversion"`
	Precheck       bool     `json:"precheck" yaml:"precheck"`     // API 호출 전 사전 체크 필요 여부
	AppContext     bool     `json:"appcontext" yaml:"appcontext"` // API 호출 전 어플리케이션 자체 Context 설정 필요 여부
	PermissionType string   `json:"permissiontype" yaml:"permissiontype"`
	Permissions    []string `json:"permissions" yaml:"permissions"`
}

type RouteGroup struct {
	Group      string  `json:"group" yaml:"group"`
	MinVersion string  `json:"minversion" yaml:"minversion"`
	MaxVersion string  `json:"maxversion" yaml:"maxversion"`
	List       []Route `json:"list" yaml:"list"`
}

// 환경설정에서 Strict로 설정된 경우 사용할 수 있는 기본 인터페이스
// IsSucceed()가 false 인 경우에 Prefix 함수는 실패한 것으로 본다.
// Prefix 함수가 실패한 경우 Response를 statusOk 와 함께 리턴하고,
// 실제 컨트롤러 함수는 실행되지 않는다.
type PreCheckResponse struct {
	IsSucceed bool
	Response  interface{}
}

type BaseController struct {
	Routes      []RouteGroup                        `json:"routes" yaml:"routes"`
	AppContext  func(echo.Context) PreCheckResponse // API 호출 전에 실행되는 어플리케이션 자체 Context 설정 함수
	PreCheck    func(echo.Context) PreCheckResponse // API 호출 전에 실행되는 사전 체크 함수
	MaxVersion  string
	Middlewares []echo.MiddlewareFunc
}

func (o *BaseController) MapRoutes(concreteController interface{}, e *echo.Echo, routefile string) error {
	if err := config.Load(routefile, o); err != nil {
		log.Errorf("Invalid routefile[%s]", routefile)
		return err
	}

	MethodRegister := map[string]func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route{
		"get":    e.GET,
		"post":   e.POST,
		"put":    e.PUT,
		"delete": e.DELETE,
	}

	for _, routeGroup := range o.Routes {
		for _, route := range routeGroup.List {
			v := reflect.ValueOf(concreteController)
			if !v.IsValid() {
				continue
			}

			f := v.MethodByName(route.Handler)
			if !f.IsValid() {
				continue
			}

			if t, ok := f.Interface().(func(echo.Context) error); ok {
				var pt func(echo.Context) error

				if route.Precheck && o.PreCheck != nil {
					log.Debug("Precheck enabled : " + route.URI)
					pt = func(permType string, perm []string) func(ec echo.Context) error {
						return func(ec echo.Context) error {
							if perm != nil && len(perm) > 0 {
								ec.Set("permissiontype", permType)
								ec.Set("permissions", perm)
							}
							res := o.PreCheck(ec)
							if !res.IsSucceed {
								return ec.JSON(http.StatusOK, res.Response)
							}
							return t(ec)
						}
					}(route.PermissionType, route.Permissions)
				} else {
					pt = t
				}

				if route.AppContext && o.AppContext != nil {
					pt = func(next func(echo.Context) error) func(ec echo.Context) error {
						return func(ec echo.Context) error {
							res := o.AppContext(ec)
							if !res.IsSucceed {
								return ec.JSON(http.StatusOK, res.Response)
							}
							return next(ec)
						}
					}(pt)
				}

				iterF := func(minVersion, maxVersion string) bool {
					if minVersion == "" || maxVersion == "" {
						return false
					}
					// if o.MaxVersion == "" || apiversion.GE(maxVersion, o.MaxVersion) {
					// 	o.MaxVersion = maxVersion
					// }

					min, err := strconv.ParseFloat(minVersion[1:], 32)
					if err != nil {
						panic("Version string may be wrong.")
					}
					max, err := strconv.ParseFloat(maxVersion[1:], 32)
					if err != nil {
						panic("Version string may be wrong.")
					}

					for ; min <= max; min += 0.1 {
						URI := strings.Replace(route.URI, ":apiver", minVersion[:1]+strconv.FormatFloat(min, 'f', 1, 64), 1)
						log.Debug("Registered: " + URI)
						MethodRegister[route.Method](URI, pt)
					}

					return true
				}

				if !iterF(route.MinVersion, route.MaxVersion) && // 각 메소드별로 버전 라우팅을 추가함(Internal API)
					!iterF(routeGroup.MinVersion, routeGroup.MaxVersion) { // 모든 메소드에 대해서 버전 라우팅을 추가(External API)
					// 호환성 유지를 위해 Min/MaxVersion 정보가 없는 라우팅테이블의 경우
					// URI 직접 매핑
					MethodRegister[route.Method](route.URI, pt)
				}
			}
		}
	}

	return nil
}
