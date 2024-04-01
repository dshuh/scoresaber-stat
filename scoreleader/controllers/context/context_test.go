package context

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/resultcode"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	webconf "stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

const (
	configFile = "../../etc/conf/config.unittest.yml"
)

var headers = []string{"content-type", "application/json"}

var host = "http://localhost:10473/"

var webConf webconf.Config
var baseContext *base.BaseContext

type args struct {
	baseCtx *base.BaseContext
}

type test struct {
	name    string
	args    args
	wantErr bool
}

func getTests(caseCount int) []test {
	var tests []test
	for i := 1; i <= caseCount; i++ {
		tests = append(tests, test{
			name: strconv.Itoa(i),
			args: args{
				baseCtx: baseContext,
			},
			wantErr: false,
		})
	}
	return tests
}

func createEcho(url, method, body string, header []string) echo.Context {

	e := echo.New()
	req := httptest.NewRequest(method, url, strings.NewReader(body))

	count := len(header)
	if count%2 != 0 {
		return nil
	}

	for i := 0; i < count; i = i + 2 {
		req.Header.Set(header[i], header[i+1])
	}

	rec := httptest.NewRecorder()
	echoContext := e.NewContext(req, rec)

	return echoContext
}

func TestMain(t *testing.T) {
	Initialize()
}

func Initialize() error {
	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	AppendRequestParameter()

	webConf = config.GetInstance(configFile).Config
	baseContext = base.NewBaseContext(createEcho(host, http.MethodPost, "", headers), &webConf)
	return nil
}

func TestNewTemplateContext(t *testing.T) {
	Initialize()

	type args struct {
		baseCtx *base.BaseContext
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "1",
			args: args{
				baseCtx: nil,
			},
			want: nil,
		},
		{
			name: "2",
			args: args{
				baseCtx: baseContext,
			},
			want: &TemplateContext{
				BaseContext: baseContext,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTemplateContext(tt.args.baseCtx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTemplateContext() = %v, want %v", got, tt.want)
			}

			if got != nil {
				ctx := got.(*TemplateContext)
				properties := []string{
					"id",
					"key",
					"value",
					"properties",
					"page_index",
					"page_size",
					"sort",
					"start_reg_at",
					"end_reg_at",
				}

				for _, v := range properties {
					ctx.SetParam(v, v)
				}

				if ctx.TemplateID() != "id" {
					t.Errorf("Property = %v, want %v", ctx.TemplateID(), "id")
				}
				if ctx.TemplateKey() != "key" {
					t.Errorf("Property = %v, want %v", ctx.UserID(), "key")
				}
				if ctx.TemplateValue() != "value" {
					t.Errorf("Property = %v, want %v", ctx.UserID(), "value")
				}
				if ctx.Properties() != "properties" {
					t.Errorf("Property = %v, want %v", ctx.Properties(), "properties")
				}
				if ctx.PageIndex() != "page_index" {
					t.Errorf("Property = %v, want %v", ctx.PageIndex(), "page_index")
				}
				if ctx.PageSize() != "page_size" {
					t.Errorf("Property = %v, want %v", ctx.PageSize(), "page_size")
				}
				if ctx.SortOrder() != "sort" {
					t.Errorf("Property = %v, want %v", ctx.SortOrder(), "sort")
				}
				if ctx.StartRegAt() != "start_reg_at" {
					t.Errorf("Property = %v, want %v", ctx.StartRegAt(), "start_reg_at")
				}
				if ctx.EndRegAt() != "end_reg_at" {
					t.Errorf("Property = %v, want %v", ctx.EndRegAt(), "end_reg_at")
				}
			}
		})
	}
}

func TestAppendRequestParameter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AppendRequestParameter()
		})
	}
}
