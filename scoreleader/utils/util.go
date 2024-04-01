package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// GetUUIDV4 UUID Version4를 얻는다. Randomic UUID > ShardKey
func GetUUIDV4() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

// GetUUIDV1 UUID Version1을 얻는다. Sequencial UUID > Index
func GetUUIDV1() string {
	return uuid.Must(uuid.NewV1(), nil).String()
}

// GetNanoTimestamp Nano Timestamp(19자리)를 얻는다.
func GetNanoTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

// GetMilliTimestamp Milli Timestamp(13자리)를 얻는다.
func GetMilliTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano()/int64(time.Millisecond), 10)
}

// GetTimestamp Timestamp(10자리)를 얻는다.
func GetTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

// HasDuplicateArray Array의 중복값을 비교한다.
func HasDuplicateArray(list, data []string) bool {
	encountered := map[string]bool{}

	for i := range data {
		encountered[data[i]] = true
	}

	for i := range list {
		if _, ok := encountered[list[i]]; ok {
			return true
		}
	}
	return false
}

// HasDuplicateString Array의 중복값을 비교한다.
func HasDuplicateString(list []string, data string) bool {
	encountered := map[string]bool{}
	encountered[data] = true

	for i := range list {
		if _, ok := encountered[list[i]]; ok {
			return true
		}
	}
	return false
}

// RemoveDuplicateArray Array의 중복값을 제거한다.
func RemoveDuplicateArray(elements []string) []string {
	encountered := map[string]bool{}

	for i := range elements {
		encountered[elements[i]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for i := range encountered {
		result = append(result, i)
	}
	return result
}

// SubtractDuplicateArray returns the elements in `a` that aren't in `b`.
func SubtractDuplicateArray(list, data []string) []string {
	mb := make(map[string]struct{}, len(data))
	for _, x := range data {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range list {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// GetPaths 계층형 지원을 위한 Path 리스트 조회
func GetPaths(path, seperator string) []string {
	if path == seperator {
		return []string{path}
	}
	pathList := strings.Split(path, seperator)
	result := make([]string, len(pathList))
	for i, v := range pathList {
		result[i] = v
		if i > 0 {
			result[i] = result[i-1] + seperator + v
		}
	}
	if result[0] == "" {
		result[0] = seperator
	}
	return result
}

// MergeStringArrayByInterfaceArray Interface Array의 요소를 String Array로 머지하고 중복값을 제거한다.
func MergeStringArrayByInterfaceArray(elements []interface{}) ([]string, error) {
	var result []string

	for _, v := range elements {
		if v == nil {
			v = "[]"
		}
		var convertStringArray []string
		if err := json.Unmarshal([]byte(v.(string)), &convertStringArray); err != nil {
			return nil, err
		}
		result = append(result, convertStringArray...)
	}
	result = RemoveDuplicateArray(result)
	return result, nil
}

// AddQueryFilter query where 구문 생성 또는 추가
func AddQueryFilter(query string, condition string) string {
	if query == "" {
		return fmt.Sprintf("where %s", condition)
	} else {
		return fmt.Sprintf("%s and %s", query, condition)
	}
}

// AddQueryFilterInt64 query where 구문 생성 또는 추가 (int64)
func AddQueryFilterInt64(query string, condition string, value int64) string {
	if value == 0 {
		return query
	}

	return AddQueryFilter(query, fmt.Sprintf("%s%d", condition, value))
}

// AddQueryFilterString query where 구문 생성 또는 추가 (string)
func AddQueryFilterString(query string, condition string, value string) string {
	if value == "" {
		return query
	}

	return AddQueryFilter(query, fmt.Sprintf("%s'%s'", condition, value))
}

// AddQueryFilterStringLike query where 구문 생성 또는 추가 (string like)
func AddQueryFilterStringLike(query string, condition string, value string) string {
	if value == "" {
		return query
	}

	return AddQueryFilter(query, fmt.Sprintf("%s like '%%%s%%'", condition, value))
}

// 현재 요청들어온 API의 버전을 조회한다.
func GetCurrentAPIVersion(path string) string {
	// 정규 표현식으로 API 버전 확인
	re := regexp.MustCompile(`^/((?:m|v)(?:1|2)\.\d+).*`)
	// re := regexp.MustCompile(`^/(?:m|v)\(?:1|2).\d+).*`)
	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// API 버전 정보를 입력받아 메이저 버전을 조회한다.
func GetCurrentMajorVersion(apiVersion string) int {
	// 정규 표현식으로 API 버전 확인
	re := regexp.MustCompile(`(?i)^[mv](\d)(?:\.\d+)?$`)
	matches := re.FindStringSubmatch(apiVersion)

	if len(matches) > 1 {
		versionInt, err := strconv.Atoi(matches[1])
		if err == nil {
			return versionInt
		}
	}
	return 0
}
