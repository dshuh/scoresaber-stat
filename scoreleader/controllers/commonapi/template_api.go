package commonapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	"stove-gitlab.sginfra.net/backend/template/models"
)

// CreateTemplate Template/01. 유저 생성
func CreateTemplate(ctx *context.TemplateContext) error {
	resp := base.Response{}

	template := models.NewTemplate()
	// RequestBody Decoding
	if err := json.NewDecoder(ctx.EchoContext.Request().Body).Decode(&template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, base.MakeError(base.ResultParameterIsInvalid, err))
	}

	// 유저 생성
	if err := models.CreateTemplate(template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, err)
	}

	resp.OK()
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// UpdateTemplate Template/02. 유저 수정
func UpdateTemplate(ctx *context.TemplateContext) error {
	resp := base.Response{}

	template := models.NewTemplate()
	// RequestBody Decoding
	if err := json.NewDecoder(ctx.EchoContext.Request().Body).Decode(&template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, base.MakeError(base.ResultParameterIsInvalid, err))
	}

	// 유저 수정
	if err := models.UpdateTemplate(template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, err)
	}

	resp.OK()
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// DeleteTemplate Template/03. 유저 삭제
func DeleteTemplate(ctx *context.TemplateContext) error {
	resp := base.Response{}

	template := models.NewTemplate()
	// RequestBody Decoding
	if err := json.NewDecoder(ctx.EchoContext.Request().Body).Decode(&template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, base.MakeError(base.ResultParameterIsInvalid, err))
	}

	// 유저 삭제
	if err := models.DeleteTemplate(template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, err)
	}

	resp.OK()
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// GetTemplate Template/04. 유저 조회
func GetTemplate(ctx *context.TemplateContext) error {
	resp := base.Response{}

	template := models.NewTemplate()
	template.Key = ctx.TemplateKey()

	// 유저 조회
	if err := models.GetTemplate(template); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, err)
	}

	resp.Value = make(map[string]interface{}, 10)

	templateConvParam := models.NewTemplateConvParam()
	templateConvParam.Properties = ctx.Properties()
	if !json.Valid([]byte(templateConvParam.Properties)) {
		templateConvParam.Properties = "[]"
	}
	templateConvParam.Data = &resp.Value

	template.ConvertTo(templateConvParam)

	resp.OK()
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// GetTemplate Template/04. 유저 조회
func GetTemplates(ctx *context.TemplateContext) error {
	resp := base.Response{}

	templateList := models.NewTemplateList()
	templateList.Key = ctx.TemplateKey()
	templateList.Value = ctx.TemplateValue()
	templateList.PageIndex, _ = strconv.Atoi(ctx.PageIndex())
	templateList.PageSize, _ = strconv.Atoi(ctx.PageSize())

	// 유저 조회
	if err := models.GetTemplates(templateList); err != nil {
		return base.JSONInternalServerError(ctx.EchoContext, err)
	}

	resp.Value = make(map[string]interface{}, 10)

	templateConvParam := models.NewTemplateConvParam()
	templateConvParam.Properties = ctx.Properties()
	if !json.Valid([]byte(templateConvParam.Properties)) {
		templateConvParam.Properties = "[]"
	}
	templateConvParam.Data = &resp.Value

	templateList.ConvertTo(templateConvParam)

	resp.Value["total_count"] = templateList.TotalCount
	resp.OK()
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
