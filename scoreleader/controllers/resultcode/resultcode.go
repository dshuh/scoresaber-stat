package resultcode

const (
	// For Template
	ResultErrorRedisConnection      = 77900
	ResultErrorAccessToken          = 77901
	ResultErrorTemplateExistName    = 77902
	ResultErrorTemplateNotFound     = 77903
	ResultErrorTemplateCreate       = 77904
	ResultErrorAccessTokenExpired   = 77905
	ResultErrorInvalidCredentials   = 77906
	ResultErrorNotAllowedClientIP   = 77907
	ResultErrorAPIParam             = 77920
	ResultErrorAPIParamKey          = 77921
	ResultErrorAPIParamValue        = 77922
	ResultErrorAPIParamDataType     = 77923
	ResultErrorBuildResultJSON      = 77924
	ResultErrorURIResource          = 77925
	ResultErrorNotFoundAPIVersion   = 77926
	ResultErrorMembershipAPI        = 77931
	ResultErrorNotFoundUserID       = 77932
	ResultErrorUserAlreadyExists    = 77933
	ResultErrorClientAlreadyExists  = 77934
	ResultErrorDataKeyAlreadyExists = 77935
	ResultErrorDataKeyExists        = 77936
	ResultErrorRegUserNotExists     = 77937
	ResultErrorClientNotExists      = 77938
	ResultErrorRequired2Characters  = 77939
	ResultErrorTemplateSystem       = 77999
)

var TemplateResultCodeMap = map[int]string{
	// For Template
	ResultErrorRedisConnection:      "error occurs when connect to DB",
	ResultErrorAccessToken:          "error occurs when user invalid access token",
	ResultErrorTemplateExistName:    "error occurs when template exist same name",
	ResultErrorTemplateNotFound:     "error occurs when template not found",
	ResultErrorTemplateCreate:       "error occurs when template create fail",
	ResultErrorAccessTokenExpired:   "access token is expired",
	ResultErrorInvalidCredentials:   "invalid user id or password",
	ResultErrorNotAllowedClientIP:   "not allowed client IP",
	ResultErrorAPIParam:             "check api parameter",
	ResultErrorAPIParamKey:          "check api parameter's default key",
	ResultErrorAPIParamValue:        "check api parameter body's data key",
	ResultErrorAPIParamDataType:     "check api parameter's data type",
	ResultErrorBuildResultJSON:      "build result json format error",
	ResultErrorURIResource:          "check uri resource name",
	ResultErrorNotFoundAPIVersion:   "not found api (check method version)",
	ResultErrorMembershipAPI:        "error occured in membership infra token api",
	ResultErrorNotFoundUserID:       "error occurs when invalid user id",
	ResultErrorUserAlreadyExists:    "user already exists",
	ResultErrorClientAlreadyExists:  "client already exists",
	ResultErrorDataKeyAlreadyExists: "data key already exists",
	ResultErrorDataKeyExists:        "data key exists",
	ResultErrorRegUserNotExists:     "reg user does not exist",
	ResultErrorClientNotExists:      "client does not exist",
	ResultErrorRequired2Characters:  "The search term must contain at least 2 characters.",
	ResultErrorTemplateSystem:       "Template System error",
}

// GetResultCodeMap Get ResultCodeMap
func GetResultCodeMap() map[int]string {
	return TemplateResultCodeMap
}
