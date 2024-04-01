package base

import (
	"fmt"
	"reflect"

	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
)

type Error string

func (e Error) Error() string { return string(e) }

type StoveError struct {
	Code    int
	Message string

	Original error
}

func (o StoveError) Error() string {
	return fmt.Sprintf("Code:%d, Message:%s, Original:%v",
		o.Code, o.Message, o.Original)
}

func MakeError(code int, err error) StoveError {
	return StoveError{Code: code, Message: ResultCodeText(code), Original: err}
}

func ToErrorString(err error) string {
	if reflect.TypeOf(err) == reflect.TypeOf(StoveError{}) {
		return err.(StoveError).Error()
	}

	return err.Error()
}

// MakeErrorWithLogf StoveError 생성하며 로그 출력을 함께 처리 합니다
func MakeErrorWithLogf(code int, format string, a ...interface{}) StoveError {
	err := MakeError(code, fmt.Errorf(format, a...))
	log.ErrorWithSkip(4, err)
	return err
}
