package base

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
)

// stove.MemberInfo에서 가져다 변형해 사용
// /amigo/amigoapp/auth/stove/membership.go
// MemberInfo 회원정보
type MemberInfo struct {
	MemberNo int    `json:"member_no"`
	Nickname string `json:"nickname"`
	Guid     string `json:"guid"`
}

type TokenResponse struct {
	Response
	Data *MemberInfo `json:"data"`
}

// Global 환경 전용 (TCloud에서 사용하면 안됨)
// token을 base64 decode해서 memberInfo를 리턴한다
func ParseToken(header http.Header) (resp TokenResponse, err error) {
	ownerContentEncoding := header.Get(GetParamField(ParamFieldOwnerContentEncoding))
	ownerContent := header.Get(GetParamField(ParamFieldOwnerContent))

	if len(ownerContentEncoding) == 0 || len(ownerContent) == 0 {
		log.Debugf("[ParseToken] invalid access token: ownerContentEncoding(%s) ownerContent(%s)", ownerContentEncoding, ownerContent)
		resp.Code = ResultInvalidAccessToken
		resp.Message = ResultCodeText(ResultInvalidAccessToken)
		return resp, nil
	}

	return parseBase64Content(ownerContent)
}

func parseBase64Content(payload string) (resp TokenResponse, err error) {
	// access token 대신 base64로 인코딩 된 정보가 header에 x-owner-content 라는 이름으로 들어온다
	// 상세내용은 http://wiki.smilegate.net:8090/pages/viewpage.action?pageId=86249038 참조

	// padding이 포함되지 않은경우 추가
	if i := len(payload) % 4; i != 0 {
		payload += strings.Repeat("=", 4-i)
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil || decoded == nil {
		log.Debugf("[ParseToken] invalid access token: payload(%s) err(%s)", payload, err.Error())
		resp.Code = ResultInvalidAccessToken
		resp.Message = ResultCodeText(ResultInvalidAccessToken)
		return resp, nil
	}

	// x-owner-content의 구조가 바뀔것을 대비해 별도로 구조체를 생성해서 unmarshal하지 않고 interface를 사용한다
	var ownerContent map[string]interface{}

	if err := json.Unmarshal([]byte(decoded), &ownerContent); err != nil {
		log.Debugf("[ParseToken] json.Unmarshal err(%s)", err.Error())
		resp.Code = ResultInvalidAccessToken
		resp.Message = ResultCodeText(ResultInvalidAccessToken)
		return resp, nil
	}

	// payload를 추출하기 위한 interface
	memberInfo := &MemberInfo{}

	// guest 회원인 경우 존재하지 않고 정회원인 경우에만 존재한다
	pld := ownerContent["pld"].(map[string]interface{})
	if pld["member_no"] != nil {
		memberInfo.MemberNo = int(pld["member_no"].(float64))
	}

	// 존재하는 경우에만 사용
	if pld["nickname"] != nil {
		memberInfo.Nickname = pld["nickname"].(string)
	}

	// add guid
	if ownerContent["guid"] != nil {
		memberInfo.Guid = ownerContent["guid"].(string)
	}

	resp.Data = memberInfo
	resp.OK()

	return resp, nil
}
