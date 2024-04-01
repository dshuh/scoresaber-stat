package base

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

// Response JSON 응답 구조체의 공통 정보를 포함한다.
type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Value   map[string]interface{} `json:"value"`
	Error   string                 `json:"error,omitempty"`
}

func ResponseInternalServerError() Response {
	return Response{
		Code:    ResultInternalServerError,
		Message: ResultCodeText(ResultInternalServerError),
	}
}

func ResponseInvalidAccessToken() Response {
	return Response{
		Code:    ResultInvalidAccessToken,
		Message: ResultCodeText(ResultInvalidAccessToken),
	}
}

func ResponseInvalidHandshakingTicket() Response {
	return Response{
		Code:    ResultInvalidHandshakingTicket,
		Message: ResultCodeText(ResultInvalidHandshakingTicket),
	}
}

func ResponseInvalidAPIVersion() Response {
	return Response{
		Code:    ResultInvalidAPIVersion,
		Message: ResultCodeText(ResultInvalidAPIVersion),
	}
}

func ResponseInvalidApplicationNo() Response {
	return Response{
		Code:    ResultInvalidApplicationNo,
		Message: ResultCodeText(ResultInvalidApplicationNo),
	}
}

func ResponseInvalidGameNo() Response {
	return Response{
		Code:    ResultInvalidGameNo,
		Message: ResultCodeText(ResultInvalidGameNo),
	}
}

func ResponseInvalidApplicationKeySecret() Response {
	return Response{
		Code:    ResultInvalidApplicationKeySecret,
		Message: ResultCodeText(ResultInvalidApplicationKeySecret),
	}
}

func ResponseDoNotHavePermission() Response {
	return Response{
		Code:    ResultDoNotHavePermissionToAccess,
		Message: ResultCodeText(ResultDoNotHavePermissionToAccess),
	}
}

func ResponseApplicationContextIsNil() Response {
	return Response{
		Code:    ResultApplicationContextIsNil,
		Message: ResultCodeText(ResultApplicationContextIsNil),
	}
}

func ResponseReadTimeout() Response {
	return Response{
		Code:    ResultReadTimeoutError,
		Message: ResultCodeText(ResultReadTimeoutError),
	}
}

func ResponseInvalidContentLength() Response {
	return Response{
		Code:    ResultContentLengthError,
		Message: ResultCodeText(ResultContentLengthError),
	}
}

func ResponseOK() Response {
	return Response{
		Code:    ResultOK,
		Message: ResultCodeText(ResultOK),
	}
}

func ResponseError(err StoveError) Response {
	originalError := ""
	if err.Original != nil {
		originalError = err.Original.Error()
	}

	return Response{
		Code:    err.Code,
		Message: err.Message,
		Error:   originalError,
	}
}

func (o *Response) OK() *Response {
	o.Code = ResultOK
	o.Message = ResultCodeText(o.Code)
	return o
}

func JSONInternalServerError(c echo.Context, err error) error {
	if reflect.TypeOf(err) != reflect.TypeOf(StoveError{}) {
		err = MakeError(ResultInternalServerError, err)
	}

	rep := ResponseError(err.(StoveError))
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONErrorCode(c echo.Context, code int) error {
	rep := ResponseError(MakeError(code, nil))
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONInvalidAPIVersion(c echo.Context) error {
	rep := ResponseInvalidAPIVersion()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONInvalidAccessToken(c echo.Context) error {
	rep := ResponseInvalidAccessToken()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONInvalidApplicationNo(c echo.Context) error {
	rep := ResponseInvalidApplicationNo()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONInvalidGameNo(c echo.Context) error {
	rep := ResponseInvalidGameNo()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONDoNotHavePermission(c echo.Context) error {
	rep := ResponseDoNotHavePermission()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONInvalidHandshakingTicket(c echo.Context) error {
	rep := ResponseInvalidHandshakingTicket()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONApplicationContextIsNil(c echo.Context) error {
	rep := ResponseApplicationContextIsNil()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONReadTimeout(c echo.Context) error {
	rep := ResponseReadTimeout()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}

func JSONInvalidContentLength(c echo.Context) error {
	rep := ResponseInternalServerError()
	c.Set(StbRetCode, rep.Code)
	return c.JSON(http.StatusOK, rep)
}
