package internalapi

import (
	"github.com/labstack/echo/v4"
	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/commonapi"
	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	webconf "stove-gitlab.sginfra.net/backend/template/framework/web/config"
	"stove-gitlab.sginfra.net/backend/template/utils"
)

// InternalAPI InternalAPI 통신을 위한 구조체
type InternalAPI struct {
	base.BaseController

	conf    *config.TemplateConfig
	apiConf *webconf.APIServer
	echo    *echo.Echo
}

// PreCheck PreCheck
func PreCheck(c echo.Context) base.PreCheckResponse {
	var result base.PreCheckResponse
	if c == nil {
		result = base.PreCheckResponse{
			IsSucceed: false,
		}
	} else {
		conf := config.GetInstance()
		if err := base.SetContext(c, &conf.Config, context.NewTemplateContext); err == nil {
			// log.Debugf("http.request.header: %v", c.Request().Header)
			ctx := base.GetContext(c).(*context.TemplateContext)

			ctx.EchoContext.Set("caller-id", conf.Template.ApplicationName)
			ctx.EchoContext.Set("caller-detail", c.Request().Header.Get("caller-detail"))

			version := utils.GetCurrentAPIVersion(c.Path())
			majorVersion := utils.GetCurrentMajorVersion(version)
			ctx.EchoContext.Set("version", version)
			ctx.EchoContext.Set("major_version", majorVersion)

			result = checkAuth(ctx)
		}
	}
	return result
}

// checkAuth 인증을 체크한다.
func checkAuth(ctx *context.TemplateContext) base.PreCheckResponse {
	// conf := config.GetInstance()
	// //cerberus := libcerberus.GetInstance()
	// cerberus := libcerberus.GetInstance(&conf.Cerberus,
	// 	nil,
	// 	conf.Template.ApplicationKey,
	// 	conf.Template.ApplicationSecret)
	// resp, _ := cerberus.AuthInternal(ctx.ApplicationKey(), ctx.ApplicationSecret())
	// if resp.Result != "000" {
	// 	return base.PreCheckResponse{
	// 		IsSucceed: false,
	// 		Response:  resp,
	// 	}
	// }

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

// Init Internal API 초기화 함수
func (o *InternalAPI) Init(e *echo.Echo) error {
	var err error
	//o.conf = config.GetInstance()
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err = o.MapRoutes(o, e, o.apiConf.Routes); err == nil {
		// serving documents for RESTful APIs
		if o.conf.Template.APIDocs {
			e.Static("/docs", "docs/int")
		}
	}
	return err
}

// GetConfig APIServer 정보를 반환한다.
func (o *InternalAPI) GetConfig() *webconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[0]
	return o.apiConf
}

// NewAPI InternalAPI 인스턴스를 생성한다.
func NewAPI() *InternalAPI {
	return &InternalAPI{}
}

// GetHealthCheck API 서버의 상태를 검사하는 용도로 사용한다.
func (o *InternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

// GetStatus 서비스 상태 체크 정보를 반환한다.
func (o *InternalAPI) GetStatus(c echo.Context) error {
	return commonapi.GetStatus(c)
}

// GetVersion 함수는 InternalAPI의 버전정보를 반환한다.
func (o *InternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}

// GetRealIP 접속한 사용자의 Real IP 정보를 반환한다.
func (o *InternalAPI) GetRealIP(c echo.Context) error {
	return commonapi.GetRealIP(c)
}

// CreateTemplate Template/01. Template 생성
func (o *InternalAPI) CreateTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.CreateTemplate(ctx)
}

// UpdateTemplate Template/02. Template 수정
func (o *InternalAPI) UpdateTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.UpdateTemplate(ctx)
}

// DeleteTemplate Template/03. Template 삭제
func (o *InternalAPI) DeleteTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.DeleteTemplate(ctx)
}

// GetTemplate Template/04. Template 조회
func (o *InternalAPI) GetTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.GetTemplate(ctx)
}

// GetTemplates Template/05. Template 목록 조회
func (o *InternalAPI) GetTemplates(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.GetTemplates(ctx)
}
