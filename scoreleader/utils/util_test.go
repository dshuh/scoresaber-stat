package utils

import (
	"strconv"
	"testing"

	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/controllers/resultcode"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
)

const (
	configFile = "../etc/conf/config.unittest.yml"
)

var conf *config.TemplateConfig

type args struct {
}

type test struct {
	name    string
	args    args
	wantErr bool
}

func getTests(caseCount int) []test {
	var tests []test
	for i := 1; i <= caseCount; i++ {
		wantErr := false
		if i == 1 {
			wantErr = true
		}
		tests = append(tests, test{
			name:    strconv.Itoa(i),
			args:    args{},
			wantErr: wantErr,
		})
	}
	return tests
}

const (
	unit_test_id = "unit_test"
)

func TestMain(t *testing.T) {
	//Initialize()
}

func Initialize() error {
	conf = config.GetInstance(configFile)

	context.AppendRequestParameter()
	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	return nil
}

func TestGetUUID(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if result := GetUUIDV4(); (result != "") != tt.wantErr {
				t.Errorf("GetUUIDV4() error = %v, wantErr %v", result, tt.wantErr)
			} else {
				t.Logf("GetUUIDV4() result = %v", result)
			}

			if result := GetUUIDV1(); (result != "") != tt.wantErr {
				t.Errorf("GetUUIDV1() error = %v, wantErr %v", result, tt.wantErr)
			} else {
				t.Logf("GetUUIDV1() result = %v", result)
			}
		})
	}
}

func TestGetTimestamp(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if result := GetNanoTimestamp(); (result != "") != tt.wantErr {
				t.Errorf("GetNanoTimestamp() error = %v, wantErr %v", result, tt.wantErr)
			} else {
				t.Logf("GetNanoTimestamp() result = %v", result)
			}

			if result := GetMilliTimestamp(); (result != "") != tt.wantErr {
				t.Errorf("GetMilliTimestamp() error = %v, wantErr %v", result, tt.wantErr)
			} else {
				t.Logf("GetMilliTimestamp() result = %v", result)
			}

			if result := GetTimestamp(); (result != "") != tt.wantErr {
				t.Errorf("GetTimestamp() error = %v, wantErr %v", result, tt.wantErr)
			} else {
				t.Logf("GetTimestamp() result = %v", result)
			}
		})
	}
}

func TestHasDuplicateArray(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := []string{"1", "2", "3"}
			data := []string{"1"}

			if result := HasDuplicateArray(list, data); result != tt.wantErr {
				t.Errorf("HasDuplicateArray() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestHasDuplicateString(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := []string{"1", "2", "3"}
			data := "1"

			if result := HasDuplicateString(list, data); result != tt.wantErr {
				t.Errorf("HasDuplicateString() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestRemoveDuplicateArray(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := []string{"1", "2", "3", "1"}

			if result := RemoveDuplicateArray(list); (len(result) > 0) != tt.wantErr {
				t.Errorf("RemoveDuplicateArray() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestSubtractDuplicateArray(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := []string{"1", "2", "3"}
			data := []string{"1"}

			if result := SubtractDuplicateArray(list, data); (len(result) > 0) != tt.wantErr {
				t.Errorf("SubtractDuplicateArray() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestGetPaths(t *testing.T) {
	Initialize()

	tests := getTests(2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path, seperator string
			tt.wantErr = true

			if tt.name == "1" {
				path = "/"
				seperator = "/"
			}
			if tt.name == "2" {
				path = "/A"
				seperator = "/"
			}

			if result := GetPaths(path, seperator); (len(result) > 0) != tt.wantErr {
				t.Errorf("GetPaths() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestMergeStringArrayByInterfaceArray(t *testing.T) {
	Initialize()

	tests := getTests(3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := make([]interface{}, 4)
			tt.wantErr = true

			if tt.name == "1" {
				elements = append(elements, "1")
				elements = append(elements, "2")
				elements = append(elements, "3")
				elements = append(elements, "4")
			}
			if tt.name == "2" {
				elements = append(elements, "1")
				elements = append(elements, "2")
				elements = append(elements, "3")
				elements = append(elements, "3")
			}
			if tt.name == "3" {
				tt.wantErr = false
				elements = append(elements, `["1"]`)
				elements = append(elements, `["2"]`)
				elements = append(elements, `["3"]`)
				elements = append(elements, `["3"]`)
			}

			if _, err := MergeStringArrayByInterfaceArray(elements); (err != nil) != tt.wantErr {
				t.Errorf("MergeStringArrayByInterfaceArray() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddQueryFilter(t *testing.T) {
	Initialize()

	tests := getTests(3)

	for _, tt := range tests {
		query := ""
		condition := ""
		if tt.name == "2" {
			tt.wantErr = true
			condition = "seq=0"
		}
		if tt.name == "3" {
			tt.wantErr = true
			query = "test"
			condition = "seq=0"
		}
		t.Run(tt.name, func(t *testing.T) {
			if result := AddQueryFilter(query, condition); (len(result) > 0) != tt.wantErr {
				t.Errorf("AddQueryFilter() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestAddQueryFilterInt64(t *testing.T) {
	Initialize()

	tests := getTests(4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := ""
			condition := ""
			value := int64(0)
			tt.wantErr = false
			if tt.name == "2" {
				condition = "seq=0"
			}
			if tt.name == "3" {
				tt.wantErr = true
				query = "test"
				condition = "seq=0"
			}
			if tt.name == "4" {
				tt.wantErr = true
				query = "test"
				condition = "seq=0"
				value = int64(1)
			}

			if result := AddQueryFilterInt64(query, condition, value); (len(result) > 0) != tt.wantErr {
				t.Errorf("AddQueryFilterInt64() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestAddQueryFilterString(t *testing.T) {
	Initialize()

	tests := getTests(4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := ""
			condition := ""
			value := ""
			tt.wantErr = false
			if tt.name == "2" {
				condition = "seq=0"
			}
			if tt.name == "3" {
				query = "test"
				condition = "seq=0"
				tt.wantErr = true
			}
			if tt.name == "3" {
				query = "test"
				condition = "seq=0"
				value = "1"
			}

			if result := AddQueryFilterString(query, condition, value); (len(result) > 0) != tt.wantErr {
				t.Errorf("AddQueryFilterString() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestAddQueryFilterStringLike(t *testing.T) {
	Initialize()

	tests := getTests(4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := ""
			condition := ""
			value := ""
			tt.wantErr = false
			if tt.name == "2" {
				condition = "seq=0"
			}
			if tt.name == "3" {
				query = "test"
				condition = "seq=0"
				tt.wantErr = true
			}
			if tt.name == "4" {
				query = "test"
				condition = "seq=0"
				value = "1"
				tt.wantErr = true
			}

			if result := AddQueryFilterStringLike(query, condition, value); (len(result) > 0) != tt.wantErr {
				t.Errorf("AddQueryFilterStringLike() result = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestGetCurrentAPIVersion(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.wantErr = false
			if result := GetCurrentAPIVersion(""); (result != "") != tt.wantErr {
				t.Errorf("GetCurrentAPIVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentAPIVersion("/v1.0/maintenance/PLATFORM/STOVE_WEB/HOME/ko"); (result != "v1.0") != tt.wantErr {
				t.Errorf("GetCurrentAPIVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentAPIVersion("/v2.0/maintenance/PLATFORM/STOVE_WEB/HOME/ko"); (result != "v2.0") != tt.wantErr {
				t.Errorf("GetCurrentAPIVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentAPIVersion("/m1.0/maintenance/PLATFORM/STOVE_WEB/HOME/ko"); (result != "m1.0") != tt.wantErr {
				t.Errorf("GetCurrentAPIVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentAPIVersion("/m2.0/maintenance/PLATFORM/STOVE_WEB/HOME/ko"); (result != "m2.0") != tt.wantErr {
				t.Errorf("GetCurrentAPIVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}

func TestGetCurrentMajorVersion(t *testing.T) {
	Initialize()

	tests := getTests(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if result := GetCurrentMajorVersion("v1.0"); (result == 1) != tt.wantErr {
				t.Errorf("GetCurrentMajorVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentMajorVersion("v2.0"); (result == 2) != tt.wantErr {
				t.Errorf("GetCurrentMajorVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentMajorVersion("m1.0"); (result == 1) != tt.wantErr {
				t.Errorf("GetCurrentMajorVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
			if result := GetCurrentMajorVersion("m2.0"); (result == 2) != tt.wantErr {
				t.Errorf("GetCurrentMajorVersion() error = %v, wantErr %v", result, tt.wantErr)
			}
		})
	}
}
