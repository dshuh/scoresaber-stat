package base

import (
	"encoding/base64"
	"strconv"
	"strings"

	"stove-gitlab.sginfra.net/backend/template/framework/web/config"

	"github.com/labstack/echo/v4"
)

const (
	maxPageSize = 1000
)

const (
	ServiceTypeMin                = iota
	ServiceTypeStoveApplicationNo // ServiceID => application_no
	ServiceTypeStoveGameNo        // ServiceID => game_no
	ServiceTypeMax
)

const (
	UserTypeMin            = iota
	UserTypeStoveMemberNo  // UserID => member_no
	UserTypeStoveProfileID // UserID => profile_id
	UserTypeMax
)

const (
	ServiceTypeStoveApplicationNoText = "application_no"
	ServiceTypeStoveGameNoText        = "game_no"
	UserTypeStoveMemberNoText         = "member_no"
	UserTypeStoveProfileIDText        = "profile_id"
)

const (
	ErrConfigIsNil        = Error("Config is nil")
	ErrInvalidUserType    = Error("Invalid user_type")
	ErrInvalidServiceType = Error("Invalid service_type")
)

const (
	ParamTypeHeader = iota
	ParamTypePath
	ParamTypeFormValue
)

const (
	ParamFieldBaseMin = iota
	ParamFieldOwnerContentType
	ParamFieldOwnerContentEncoding
	ParamFieldOwnerContent
	ParamFieldAuthorization
	ParamFieldAccessToken
	ParamFieldApplicationKey
	ParamFieldApplicationSecret
	ParamFieldServerApplicationKey
	ParamFieldServerApplicationSecret
	ParamFieldClientApplicationKey
	ParamFieldClientApplicationSecret
	ParamFieldServerApplicationNo
	ParamFieldApplicationName
	ParamFieldApplicationDescription
	ParamFieldExistsApplicationNo
	ParamFieldGameNo
	ParamFieldDuplicateYN
	ParamFieldPageOffset
	ParamFieldPageSize
	ParamFieldBaseMax
)

type ParamType int
type ParamField int
type ParamMeta struct {
	ParamType ParamType
	FieldName string
}

var serviceTypeText = map[string]int{
	ServiceTypeStoveApplicationNoText: ServiceTypeStoveApplicationNo,
	ServiceTypeStoveGameNoText:        ServiceTypeStoveGameNo,
}

var userTypeText = map[string]int{
	UserTypeStoveMemberNoText:  UserTypeStoveMemberNo,
	UserTypeStoveProfileIDText: UserTypeStoveProfileID,
}

// RequestParameter Cerberus API를 요청할 때 사용하는 입력 매개변수의 메타데이터
var requestParameter = map[ParamField]*ParamMeta{
	ParamFieldOwnerContentType:        NewParamMeta(ParamTypeHeader, "x-owner-content-type"),
	ParamFieldOwnerContentEncoding:    NewParamMeta(ParamTypeHeader, "x-owner-content-encoding"),
	ParamFieldOwnerContent:            NewParamMeta(ParamTypeHeader, "x-owner-content"),
	ParamFieldAuthorization:           NewParamMeta(ParamTypeHeader, "Authorization"),
	ParamFieldAccessToken:             NewParamMeta(ParamTypeHeader, "Access-Token"),
	ParamFieldApplicationKey:          NewParamMeta(ParamTypeHeader, "Application-Key"),
	ParamFieldApplicationSecret:       NewParamMeta(ParamTypeHeader, "Application-Secret"),
	ParamFieldServerApplicationKey:    NewParamMeta(ParamTypeHeader, "Server-Application-Key"),
	ParamFieldServerApplicationSecret: NewParamMeta(ParamTypeHeader, "Server-Application-Secret"),
	ParamFieldClientApplicationKey:    NewParamMeta(ParamTypeHeader, "Client-Application-Key"),
	ParamFieldClientApplicationSecret: NewParamMeta(ParamTypeHeader, "Client-Application-Secret"),
	ParamFieldServerApplicationNo:     NewParamMeta(ParamTypePath, "server_application_no"),
	ParamFieldApplicationName:         NewParamMeta(ParamTypeFormValue, "application_name"),
	ParamFieldApplicationDescription:  NewParamMeta(ParamTypeFormValue, "application_description"),
	ParamFieldExistsApplicationNo:     NewParamMeta(ParamTypeFormValue, "exists_application_no"),
	ParamFieldGameNo:                  NewParamMeta(ParamTypeFormValue, "game_no"),
	ParamFieldDuplicateYN:             NewParamMeta(ParamTypeFormValue, "duplicate_yn"),
	ParamFieldPageOffset:              NewParamMeta(ParamTypeFormValue, "page_offset"),
	ParamFieldPageSize:                NewParamMeta(ParamTypeFormValue, "page_size"),
}

func NewParamMeta(paramType ParamType, fieldName string) *ParamMeta {
	return &ParamMeta{ParamType: paramType, FieldName: fieldName}
}

func GetParamField(field ParamField) string {
	return requestParameter[field].FieldName
}

func SetParamMeta(field ParamField, paramMeta *ParamMeta) {
	requestParameter[field] = paramMeta
}

type BaseContext struct {
	EchoContext echo.Context

	conf   *config.Config
	params map[string]string

	Error error
}

func NewBaseContext(echoContext echo.Context, conf *config.Config) *BaseContext {
	ctx := new(BaseContext)
	ctx.EchoContext = echoContext
	ctx.conf = conf
	ctx.params = make(map[string]string)

	return ctx
}

func GetContext(c echo.Context) interface{} {
	ctx := c.Get("appcontext")

	return ctx
}

// NOTE: SetContext()에서 parseParams()를 할 경우 불필요한 연산이 발생되므로,
// ParseParams()를 주석처리하고 Export해서 Application에서 선택적으로 처리할 수 있도록 개선 필요
func SetContext(c echo.Context, conf *config.Config, newCtx func(*BaseContext) interface{}) error {
	if conf == nil {
		return ErrConfigIsNil
	}

	ctx := NewBaseContext(c, conf)
	ctx.parseParams()

	appCtx := newCtx(ctx)
	c.Set("appcontext", appCtx)

	return nil
}

func (o *BaseContext) SetParam(key string, value string) {
	o.params[key] = value
}

func (o *BaseContext) GetParam(field ParamField) string {
	return o.params[GetParamField(field)]
}

func (o *BaseContext) parseParams() {
	paramTypeGetFuncs := make(map[ParamType]func(string) string)

	paramTypeGetFuncs[ParamTypeHeader] = o.EchoContext.Request().Header.Get
	paramTypeGetFuncs[ParamTypePath] = o.EchoContext.Param
	paramTypeGetFuncs[ParamTypeFormValue] = o.EchoContext.FormValue
	for _, param := range requestParameter {
		o.params[param.FieldName] = paramTypeGetFuncs[param.ParamType](param.FieldName)
	}

	o.parseAuthorization()
}

func (o *BaseContext) parseAuthorization() {
	authorization, ok := o.params[requestParameter[ParamFieldAuthorization].FieldName]

	if !ok || authorization == "" {
		return
	}

	authTokens := strings.Split(authorization, " ")
	if len(authTokens) != 2 {
		return
	}

	switch strings.ToLower(authTokens[0]) {
	case "basic":
		decoded, err := base64.StdEncoding.DecodeString(authTokens[1])
		if err != nil {
			return
		}
		appTokens := strings.Split(string(decoded[:]), ":")
		if len(appTokens) != 2 {
			return
		}
		o.SetParam(GetParamField(ParamFieldApplicationKey), appTokens[0])
		o.SetParam(GetParamField(ParamFieldApplicationSecret), appTokens[1])
	case "bearer":
		o.SetParam(GetParamField(ParamFieldAccessToken), authTokens[1])
	}
}

func (o BaseContext) UserID() string {
	return o.EchoContext.FormValue("user_id")
}

func (o BaseContext) UserType() string {
	userType := o.EchoContext.FormValue("user_type")

	_, ok := userTypeText[userType]
	if ok {
		return userType
	}

	// 기본 타입 반환
	return UserTypeStoveMemberNoText
}

func (o BaseContext) ServiceID() string {
	return o.EchoContext.FormValue("service_id")
}

func (o BaseContext) ServiceType() string {
	serviceType := o.EchoContext.FormValue("service_type")

	_, ok := serviceTypeText[serviceType]
	if ok {
		return serviceType
	}

	// 기본 타입 반환
	return ServiceTypeStoveApplicationNoText
}

func (o *BaseContext) AccessToken() string {
	return o.params[requestParameter[ParamFieldAccessToken].FieldName]
}

// OwnerContentType 인증 컨텐츠 타입
func (o *BaseContext) OwnerContentType() string {
	return o.GetParam(ParamFieldOwnerContentType)
}

// OwnerContentEncoding 인증 인코딩 타입
func (o *BaseContext) OwnerContentEncoding() string {
	return o.GetParam(ParamFieldOwnerContentEncoding)
}

// OwnerContent 인증 컨텐츠
func (o *BaseContext) OwnerContent() string {
	return o.GetParam(ParamFieldOwnerContent)
}

func MakeApplicationID(key, secret string) string {
	return key + ":" + secret
}

func MakeBasicAuthorization(key, secret string) string {
	return "basic " + base64.StdEncoding.EncodeToString([]byte(MakeApplicationID(key, secret)))
}

func MakeBearerAuthorization(accessToken string) string {
	return "bearer " + accessToken
}

// ApplicationID key와 secret의 조합으로 ID 생성
// NOTE: key와 secret으로 hash 문자열을 만들 경우 비용이 크게 증가하여, 단순 합으로 반환
func (o *BaseContext) ApplicationID() string {
	return MakeApplicationID(o.ApplicationKey(), o.ApplicationSecret())
}

func (o *BaseContext) ApplicationKey() string {
	return o.params[requestParameter[ParamFieldApplicationKey].FieldName]
}

func (o *BaseContext) ApplicationSecret() string {
	return o.params[requestParameter[ParamFieldApplicationSecret].FieldName]
}

func (o *BaseContext) ServerApplicationID() string {
	return MakeApplicationID(o.ServerApplicationKey(), o.ServerApplicationSecret())
}

func (o *BaseContext) ClientApplicationID() string {
	return MakeApplicationID(o.ClientApplicationKey(), o.ClientApplicationSecret())
}

func (o *BaseContext) ServerApplicationKey() string {
	return o.GetParam(ParamFieldServerApplicationKey)
}

func (o *BaseContext) ServerApplicationSecret() string {
	return o.GetParam(ParamFieldServerApplicationSecret)
}

func (o *BaseContext) ClientApplicationKey() string {
	return o.GetParam(ParamFieldClientApplicationKey)
}

func (o *BaseContext) ClientApplicationSecret() string {
	return o.GetParam(ParamFieldClientApplicationSecret)
}

func (o *BaseContext) ServerApplicationNo() string {
	return o.GetParam(ParamFieldServerApplicationNo)
}

func (o *BaseContext) ApplicationName() string {
	return o.GetParam(ParamFieldApplicationName)
}

func (o *BaseContext) ApplicationDescription() string {
	return o.GetParam(ParamFieldApplicationDescription)
}

func (o *BaseContext) ExistsApplicationNo() string {
	return o.GetParam(ParamFieldExistsApplicationNo)
}

func (o *BaseContext) GameNo() string {
	return o.GetParam(ParamFieldGameNo)
}

func (o *BaseContext) DuplicateYN() string {
	return o.GetParam(ParamFieldDuplicateYN)
}

func (o *BaseContext) PageOffset() string {
	return o.GetParam(ParamFieldPageOffset)
}

func (o *BaseContext) PageOffsetInt64() (int64, error) {
	return strconv.ParseInt(o.GetParam(ParamFieldPageOffset), 10, 64)
}

func (o *BaseContext) PageSize() string {
	return o.GetParam(ParamFieldPageSize)
}

func (o *BaseContext) PageSizeInt() (int, error) {
	i, err := strconv.Atoi(o.GetParam(ParamFieldPageSize))

	if err != nil {
		return i, err
	} else if i > maxPageSize {
		return maxPageSize, nil
	}
	return i, nil
}
