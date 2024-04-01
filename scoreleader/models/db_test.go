package models

import (
	"strconv"
	"testing"

	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/controllers/resultcode"
	"stove-gitlab.sginfra.net/backend/template/framework/database/mysqldb"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	"stove-gitlab.sginfra.net/backend/template/utils"

	"stove-gitlab.sginfra.net/backend/template/config"
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
	unit_test_id  = "unit_test"
	unit_test_id2 = "unit_test2"
	unit_test_id3 = "unit_test2"
)

func GetUnitTestTmp(tmp string) string {
	return unit_test_id + "_" + tmp
}

func TestMain(t *testing.T) {
	//Initialize()
}

func Initialize() error {
	conf = config.GetInstance(configFile)

	context.AppendRequestParameter()
	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	return nil
}

func Test_CreateTable(t *testing.T) {
	Initialize()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "CreateTable",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateTable(); (err != nil) != tt.wantErr {
				t.Errorf("CreateTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_GetDBStatus(t *testing.T) {
	Initialize()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "GetDBStatus",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetDBStatus(); (err != nil) != tt.wantErr {
				t.Errorf("GetDBStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateServiceDB_CreateTemplate(t *testing.T) {
	Initialize()

	type args struct {
		template *Template
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				template: &Template{
					Key:   unit_test_id,
					Value: unit_test_id,
					RegAt: utils.GetTimestamp(),
					UpdAt: utils.GetTimestamp(),
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				template: &Template{
					Key:   unit_test_id2,
					Value: unit_test_id2,
					RegAt: utils.GetTimestamp(),
					UpdAt: utils.GetTimestamp(),
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				template: &Template{
					Key:   unit_test_id3,
					Value: unit_test_id3,
					RegAt: utils.GetTimestamp(),
					UpdAt: utils.GetTimestamp(),
				},
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				template: &Template{
					Key:   "",
					Value: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateTemplate(tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateServiceDB_UpdateTemplate(t *testing.T) {
	Initialize()

	type args struct {
		template *Template
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				template: &Template{
					Key:   unit_test_id,
					Value: unit_test_id2,
					UpdAt: utils.GetTimestamp(),
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				template: &Template{
					Key:   unit_test_id2,
					Value: unit_test_id,
					UpdAt: utils.GetTimestamp(),
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				template: &Template{
					Key:   "",
					Value: "",
					UpdAt: utils.GetTimestamp(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateTemplate(tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateServiceDB_GetTemplate(t *testing.T) {
	Initialize()

	type args struct {
		template *Template
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				template: &Template{
					Key: unit_test_id,
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				template: &Template{
					Key: unit_test_id2,
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				template: &Template{
					Key: unit_test_id3,
				},
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				template: &Template{
					Key: "sdfasdf",
				},
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				template: &Template{
					Key: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "3" {
				tmpTemplate := NewTemplate()
				tmpTemplate.Key = unit_test_id3
				RemoveTemplateRedis(tmpTemplate)
			}
			if err := GetTemplate(tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "3" {
				tmpTemplate := NewTemplate()
				tmpTemplate.Key = unit_test_id3
				tmpTemplate.Value = unit_test_id3
				SetTemplateRedis(tmpTemplate)
			}
		})
	}
}

func TestTemplateServiceDB_GetTemplates(t *testing.T) {
	Initialize()

	type args struct {
		templateList *TemplateList
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				templateList: &TemplateList{
					Key:       unit_test_id,
					PageIndex: 1,
					PageSize:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				templateList: &TemplateList{
					Value:     unit_test_id,
					PageIndex: 1,
					PageSize:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				templateList: &TemplateList{
					Key:       unit_test_id2,
					PageIndex: 1,
					PageSize:  1,
					SortOrder: "desc",
				},
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				templateList: &TemplateList{
					Value:     unit_test_id2,
					PageIndex: 1,
					PageSize:  1,
					SortOrder: "asc",
				},
			},
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				templateList: &TemplateList{
					Key:   "",
					Value: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetTemplates(tt.args.templateList); (err != nil) != tt.wantErr {
				t.Errorf("GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}

			data := make(map[string]interface{}, 10)
			templateConvParam := NewTemplateConvParam()
			if tt.name == "1" {
				templateConvParam.Properties = `[sss"id","key","value", "reg_at", "upd_at"]`
			} else {
				templateConvParam.Properties = `["id","key","value", "reg_at", "upd_at"]`
			}
			templateConvParam.Data = &data

			tt.args.templateList.ConvertTo(templateConvParam)
		})
	}
}

func TestTemplateServiceDB_RefreshRedis(t *testing.T) {
	Initialize()

	type args struct {
		templateList *TemplateList
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				templateList: &TemplateList{
					Key:       unit_test_id,
					PageIndex: 1,
					PageSize:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				templateList: &TemplateList{
					Value:     unit_test_id,
					PageIndex: 1,
					PageSize:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				templateList: &TemplateList{
					Key:       unit_test_id2,
					PageIndex: 1,
					PageSize:  1,
					SortOrder: "desc",
				},
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				templateList: &TemplateList{
					Value:     unit_test_id2,
					PageIndex: 1,
					PageSize:  1,
					SortOrder: "asc",
				},
			},
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				templateList: &TemplateList{
					PageIndex: 1,
					PageSize:  1000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RefreshRedis(tt.args.templateList); (err != nil) != tt.wantErr {
				t.Errorf("RefreshRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getTable(t *testing.T) {
	Initialize()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "getTable",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if get := getTemplateTable(); get == "" {
				t.Errorf("getTemplateTable() error = %v", get)
			}
		})
	}
}

func CreateTemplateTestData() {
	CreateTemplate(&Template{Key: unit_test_id, Value: unit_test_id})
}

func DeleteTemplateTestData() {
	DeleteTemplate(&Template{Key: unit_test_id})
	DeleteTemplate(&Template{Key: unit_test_id2})
}

func TestTemplateServiceDB_CreateTemplateTable(t *testing.T) {
	Initialize()

	type fields struct {
		conf    *config.TemplateConfig
		mysqldb *mysqldb.Mysql
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "CreateTemplateTable",
			fields: fields{
				conf: config.GetInstance(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &TemplateServiceDB{
				conf:    tt.fields.conf,
				mysqldb: tt.fields.mysqldb,
			}
			if err := o.CreateTemplateTable(); (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplateTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateServiceDB_DeleteTemplate(t *testing.T) {
	Initialize()

	type args struct {
		template *Template
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				template: &Template{
					Key: unit_test_id,
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				template: &Template{
					Key: unit_test_id2,
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				template: &Template{
					Key: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteTemplate(tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateServiceDB_ValidTemplateParams(t *testing.T) {
	Initialize()

	type args struct {
		ID    int
		Key   string
		Value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ID:    1,
				Key:   "test",
				Value: "test",
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				ID:    0,
				Key:   "",
				Value: "",
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				ID:    1,
				Key:   "",
				Value: "",
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ID:    1,
				Key:   "test",
				Value: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := NewTemplate()
			template.ID = tt.args.ID
			template.Key = tt.args.Key
			template.Value = tt.args.Value

			if err := validTemplateParams([]string{"id", "key", "value"}, template); (err != nil) != tt.wantErr {
				t.Errorf("validTemplateParams error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateServiceDB_TemplateConvertTo(t *testing.T) {
	Initialize()

	type args struct {
		ID         int
		Key        string
		Value      string
		Properties string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ID:         1,
				Key:        "test",
				Value:      "test",
				Properties: `["key_id","key_name","reg_at","upd_at"]`,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				ID:         1,
				Key:        "test",
				Value:      "test",
				Properties: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := NewTemplate()
			template.ID = tt.args.ID
			template.Key = tt.args.Key
			template.Value = tt.args.Value

			templateList := NewTemplateList()
			templateList.Templates = append(templateList.Templates, *template)
			templateList.TotalCount = 1

			var data = make(map[string]interface{}, 10)
			templateConvParam := NewTemplateConvParam()
			templateConvParam.Properties = tt.args.Properties
			templateConvParam.Data = &data

			if err := template.ConvertTo(templateConvParam); (err != nil) != tt.wantErr {
				t.Errorf("Template.ConvertTo() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := templateList.ConvertTo(templateConvParam); (err != nil) != tt.wantErr {
				t.Errorf("TemplateList.ConvertTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_InitTable(t *testing.T) {
	Initialize()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "InitTable",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitTable(); (err != nil) != tt.wantErr {
				t.Errorf("InitTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_PingRedis(t *testing.T) {
	Initialize()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "PingRedis",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := PingRedis(); (err != nil) != tt.wantErr {
				t.Errorf("PingRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
