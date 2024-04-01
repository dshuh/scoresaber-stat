package base

const (
	// for http
	ResultOK                  = 0
	ResultSuccess             = 200
	ResultBadRequest          = 400
	ResultInternalServerError = 500
	ResultReadTimeoutError    = 600
	ResultContentLengthError  = 700

	// for stove
	ResultInvalidAccessToken          = 70001
	ResultInvalidHandshakingTicket    = 70002
	ResultDBIsNotInitialized          = 70003
	ResultInvalidApplicationKeySecret = 70004
	ResultStorageIsNotInitialized     = 70005
	ResultInvalidAPIVersion           = 70006
	ResultInvalidApplicationNo        = 70007
	ResultInvalidGameNo               = 70008
	ResultDoNotHavePermissionToAccess = 70009
	ResultApplicationContextIsNil     = 70010
	ResultUserIDIsRequired            = 70021
	ResultUserTypeIsRequired          = 70022
	ResultServiceIDIsRequired         = 70023
	ResultServiceTypeIsRequired       = 70024
	ResultParameterIsInvalid          = 70051
	ResultCannotFindData              = 70052
)

var resultCodeText = map[int]string{
	// for http
	ResultOK:                  "OK",
	ResultSuccess:             "success",
	ResultBadRequest:          "Bad request",
	ResultInternalServerError: "Internal server error",
	ResultReadTimeoutError:    "HTTP read timeout error",
	ResultContentLengthError:  "Content-Length is not equal payload size",

	// for stove
	ResultInvalidAccessToken:          "Invalid access token",
	ResultInvalidHandshakingTicket:    "Invalid Handshaking ticket",
	ResultDBIsNotInitialized:          "DB is not initialized",
	ResultInvalidApplicationKeySecret: "Invalid application key/secret",
	ResultStorageIsNotInitialized:     "Storage is not initialized",
	ResultInvalidAPIVersion:           "Invalid API version",
	ResultInvalidApplicationNo:        "Invalid Application No",
	ResultInvalidGameNo:               "Invalid Game No",
	ResultDoNotHavePermissionToAccess: "You do not have permission to access the resource",
	ResultApplicationContextIsNil:     "Application context is nil",
	ResultUserIDIsRequired:            "user_id is required",
	ResultUserTypeIsRequired:          "user_type is required",
	ResultServiceIDIsRequired:         "service_id is required",
	ResultServiceTypeIsRequired:       "service_type is required",
	ResultParameterIsInvalid:          "Parameter is invalid",
	ResultCannotFindData:              "Cannot find data",
}

func ResultCodeText(code int) string {
	return resultCodeText[code]
}

//자신의 프로젝트에서 정의한 ResultCodeText 데이터를 기존의 데이터에 추가한다
//주의! 중복 키의 경우 기존의 데이터가 갱신됨. 공통 오류코드 영역은 새로 정의하면 안된다
func AppendResultCodeText(appendMap *map[int]string) error {
	for k, v := range *appendMap {
		resultCodeText[k] = v
	}
	return nil
}
