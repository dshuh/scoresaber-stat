package externalapi

import (
	"net/http"

	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	webconf "stove-gitlab.sginfra.net/backend/template/framework/web/config"

	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/commonapi"
	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/utils"

	"github.com/labstack/echo/v4"
)

// ExternalAPI External 통신을 위한 구조체
type ExternalAPI struct {
	base.BaseController

	conf    *config.TemplateConfig
	echo    *echo.Echo
	apiConf *webconf.APIServer
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
			ctx := base.GetContext(c).(*context.TemplateContext)

			method := c.Request().Method
			path := c.Request().URL.Path
			header := c.Request().Header

			ctx.EchoContext.Set("caller-id", conf.Template.ApplicationName)
			ctx.EchoContext.Set("caller-detail", header.Get("caller-detail"))

			version := utils.GetCurrentAPIVersion(c.Path())
			majorVersion := utils.GetCurrentMajorVersion(version)
			ctx.EchoContext.Set("version", version)
			ctx.EchoContext.Set("major_version", majorVersion)

			result = checkAuth(ctx, method, path, header)
		}
		// base.SetContext(c, &conf.Config, context.NewTemplateContext)
		// result = base.PreCheckResponse{
		// 	IsSucceed: true,
		// }
	}
	return result
}

// checkAuth 인증을 체크한다.
func checkAuth(ctx *context.TemplateContext, method, path string, header http.Header) base.PreCheckResponse {
	var resp base.TokenResponse
	var err error

	if resp, err = base.ParseToken(header); err != nil {
		return base.PreCheckResponse{
			IsSucceed: false,
			Response:  resp,
		}
	}

	if resp.Code == base.ResultOK {
		ctx.EchoContext.Set("member_no", resp.Data.MemberNo)
	}

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

// Init External API 초기화 함수
func (o *ExternalAPI) Init(e *echo.Echo) error {
	var err error
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err == nil {
		// serving documents for RESTful APIs
		if o.conf.Template.APIDocs {
			e.Static("/docs", "docs/ext")
		}
	}
	return err
}

// GetConfig APIServer 정보를 반환한다.
func (o *ExternalAPI) GetConfig() *webconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[1]
	return o.apiConf
}

// NewAPI ExternalAPI 인스턴스를 생성한다.
func NewAPI() *ExternalAPI {
	return &ExternalAPI{}
}

// GetHealthCheck API 서버의 상태를 검사하는 용도로 사용한다.
func (o *ExternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

// GetVersion 함수는 ExternalAPI의 버전정보를 반환한다.
func (o *ExternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}

// CreateTemplate Template/01. Template 생성
func (o *ExternalAPI) CreateTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.CreateTemplate(ctx)
}

// UpdateTemplate Template/02. Template 수정
func (o *ExternalAPI) UpdateTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.UpdateTemplate(ctx)
}

// DeleteTemplate Template/03. Template 삭제
func (o *ExternalAPI) DeleteTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.DeleteTemplate(ctx)
}

// GetTemplate Template/04. Template 조회
func (o *ExternalAPI) GetTemplate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.GetTemplate(ctx)
}

// GetTemplates Template/05. Template 목록 조회
func (o *ExternalAPI) GetTemplates(c echo.Context) error {
	ctx := base.GetContext(c).(*context.TemplateContext)
	return commonapi.GetTemplates(ctx)
}
