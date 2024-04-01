package models

import (
	"encoding/json"
	"errors"

	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
)

// Template Template 정보
type Template struct {
	ID    int    `json:"id"`               // ID
	Key   string `json:"key,omitempty"`    // KeyType
	Value string `json:"value,omitempty"`  // KeyValue
	RegAt string `json:"reg_at,omitempty"` // RegAt
	UpdAt string `json:"upd_at,omitempty"` // UpdAt
}

// TemplateList TemplateList 요청/응답 정보
type TemplateList struct {
	// 요청 정보
	Key        string `json:"key"`          // Key
	Value      string `json:"value"`        // Value
	PageIndex  int    `json:"page_index"`   // 페이지 번호
	PageSize   int    `json:"page_size"`    // 한 페이지에 보여줄 갯수
	SortOrder  string `json:"sort"`         // 정렬 방식 asc, desc
	StartRegAt string `json:"start_reg_at"` // 생성일 시작 시간
	EndRegAt   string `json:"end_reg_at"`   // 생성일 종료 시간

	Properties []string `json:"properties"` // 조회할 필드 리스트

	// 결과 값
	Templates  []Template `json:"templates"`
	TotalCount int        `json:"total_count"` // 요청 정보의 총 데이터 갯수.
}

// TemplateConvParam - Template 데이터 변환을 위한 파라미터 정보를 가진 구조체
type TemplateConvParam struct {
	Templates  *[]Template
	Properties string
	Data       *map[string]interface{}
	Overwrite  bool
	Props      *[]string
}

// NewTemplate Template 참조를 반환한다.
func NewTemplate() *Template {
	return &Template{}
}

// NewTemplateList TemplateList 참조를 반환한다.
func NewTemplateList() *TemplateList {
	return &TemplateList{}
}

// NewTemplateConvParam TemplateConvParam 생성 및 초기화
func NewTemplateConvParam() *TemplateConvParam {
	return &TemplateConvParam{}
}

// ConvertTemplateListTo Template List로 변환
func ConvertTemplateListTo(templateConvParam *TemplateConvParam) error {
	if err := json.Unmarshal([]byte(templateConvParam.Properties), &templateConvParam.Props); err != nil {
		return err
	}

	list := make([]map[string]interface{}, len(*templateConvParam.Templates))

	for i, v := range *templateConvParam.Templates {
		list[i] = make(map[string]interface{}, len(*templateConvParam.Props))

		tmpTemplateParam := NewTemplateConvParam()
		tmpTemplateParam.Data = &list[i]
		tmpTemplateParam.Props = templateConvParam.Props
		tmpTemplateParam.Overwrite = templateConvParam.Overwrite

		v.ConvertToWithProps(tmpTemplateParam)
	}
	(*templateConvParam.Data)["list"] = list
	return nil
}

// ConvertTo Template로 변환
func (o *Template) ConvertTo(templateConvParam *TemplateConvParam) error {
	if err := json.Unmarshal([]byte(templateConvParam.Properties), &templateConvParam.Props); err != nil {
		return err
	}
	return o.ConvertToWithProps(templateConvParam)
}

// ConvertTo TemplateList로 변환
func (o *TemplateList) ConvertTo(templateConvParam *TemplateConvParam) error {
	templateConvParam.Templates = &o.Templates
	return ConvertTemplateListTo(templateConvParam)
}

// ConvertToWithProps Template로 변환
func (o *Template) ConvertToWithProps(templateConvParam *TemplateConvParam) error {
	for _, v := range *templateConvParam.Props {
		if _, ok := (*templateConvParam.Data)[v]; ok {
			continue
		}

		switch v {
		case "id":
			(*templateConvParam.Data)[v] = o.ID
		case "key":
			(*templateConvParam.Data)[v] = o.Key
		case "value":
			(*templateConvParam.Data)[v] = o.Value
		case "reg_at":
			(*templateConvParam.Data)[v] = o.RegAt
		case "upd_at":
			(*templateConvParam.Data)[v] = o.UpdAt
		}
	}

	return nil
}

func validTemplateParams(param []string, template *Template) error {
	for i := 0; i < len(param); i++ {
		if err := validTemplateParam(param[i], template); err != nil {
			return err
		}
	}
	return nil
}

func validTemplateParam(param string, template *Template) error {
	switch param {
	case "id":
		if template.ID == 0 {
			return base.MakeError(base.ResultParameterIsInvalid, errors.New(param))
		}
	case "key":
		if template.Key == "" {
			return base.MakeError(base.ResultParameterIsInvalid, errors.New(param))
		}
	case "value":
		if template.Value == "" {
			return base.MakeError(base.ResultParameterIsInvalid, errors.New(param))
		}
	}
	return nil
}
