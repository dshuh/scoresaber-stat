package context

import (
	"github.com/labstack/echo/v4"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
)

const (
	// ParamFieldID id 입력 파라미터 Key
	TemplateID = iota + base.ParamFieldBaseMax + 1
	TemplateKey
	TemplateValue
	Properties
	PageIndex
	PageSize
	SortOrder
	StartRegAt
	EndRegAt
)

// TemplateContext Template API의 Request Context
type TemplateContext struct {
	*base.BaseContext
	userID   string
	userName string
}

// NewTemplateContext 새로운 Template Context 생성
func NewTemplateContext(baseCtx *base.BaseContext) interface{} {
	if baseCtx == nil {
		return nil
	}

	ctx := new(TemplateContext)
	ctx.BaseContext = baseCtx

	return ctx
}

// AppendRequestParameter BaseContext에 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
	base.SetParamMeta(TemplateID, base.NewParamMeta(base.ParamTypeFormValue, "id"))
	base.SetParamMeta(TemplateKey, base.NewParamMeta(base.ParamTypeFormValue, "key"))
	base.SetParamMeta(TemplateValue, base.NewParamMeta(base.ParamTypeFormValue, "value"))
	base.SetParamMeta(Properties, base.NewParamMeta(base.ParamTypeFormValue, "properties"))
	base.SetParamMeta(PageIndex, base.NewParamMeta(base.ParamTypeFormValue, "page_index"))
	base.SetParamMeta(PageSize, base.NewParamMeta(base.ParamTypeFormValue, "page_size"))
	base.SetParamMeta(SortOrder, base.NewParamMeta(base.ParamTypeFormValue, "sort"))
	base.SetParamMeta(StartRegAt, base.NewParamMeta(base.ParamTypeFormValue, "start_reg_at"))
	base.SetParamMeta(EndRegAt, base.NewParamMeta(base.ParamTypeFormValue, "end_reg_at"))
}

func (o TemplateContext) TemplateID() string {
	return o.GetParam(TemplateID)
}

func (o TemplateContext) TemplateKey() string {
	return o.GetParam(TemplateKey)
}

func (o TemplateContext) TemplateValue() string {
	return o.GetParam(TemplateValue)
}
func (o TemplateContext) Properties() string {
	return o.GetParam(Properties)
}

func (o TemplateContext) PageIndex() string {
	return o.GetParam(PageIndex)
}

func (o TemplateContext) PageSize() string {
	return o.GetParam(PageSize)
}

func (o TemplateContext) SortOrder() string {
	return o.GetParam(SortOrder)
}

func (o TemplateContext) StartRegAt() string {
	return o.GetParam(StartRegAt)
}

func (o TemplateContext) EndRegAt() string {
	return o.GetParam(EndRegAt)
}

func (o TemplateContext) GetRealIP() string {
	return o.BaseContext.EchoContext.Request().Header.Get(echo.HeaderXRealIP)
}
