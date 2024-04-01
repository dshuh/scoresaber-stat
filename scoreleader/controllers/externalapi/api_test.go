package externalapi

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/controllers/resultcode"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	webconf "stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

const (
	configFile = "../../etc/conf/config.unittest.yml"
)

var headers = []string{"content-type", "application/json"}

var host = "http://localhost:10473/"

var extAPI *ExternalAPI
var conf *config.TemplateConfig
var echoContext echo.Context
var templateCtx *context.TemplateContext

func TestMain(t *testing.T) {
	Initialize()
}

func Initialize() error {
	conf = config.GetInstance(configFile)

	extAPI = NewAPI()
	extAPI.apiConf = extAPI.GetConfig()
	extAPI.echo = echo.New()

	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	context.AppendRequestParameter()
	echoContext = createEcho(host, http.MethodGet, "", headers)

	templateCtx = base.GetContext(echoContext).(*context.TemplateContext)
	return nil
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
	echoContext = e.NewContext(req, rec)

	if err := base.SetContext(echoContext, &conf.Config, context.NewTemplateContext); err != nil {
		log.Error(err)
	}

	return echoContext
}

func TestExternalAPI_Init(t *testing.T) {
	//Initialize()
	if extAPI.apiConf.Routes[0:1] != "/" {
		extAPI.apiConf.Routes = "../../etc/conf/" + extAPI.apiConf.Routes
	}

	type args struct {
		e *echo.Echo
	}
	tests := []struct {
		name    string
		o       *ExternalAPI
		args    args
		wantErr bool
	}{
		{
			name: "ExternalAPI_Init",
			o:    extAPI,
			args: args{
				e: extAPI.echo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Init(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreCheck(t *testing.T) {
	//Initialize()
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name string
		args args
		want base.PreCheckResponse
	}{
		{
			name: "PreCheck",
			args: args{
				c: echoContext,
			},
			want: base.PreCheckResponse{
				IsSucceed: true,
			},
		},
		{
			name: "PreCheck",
			args: args{
				c: nil,
			},
			want: base.PreCheckResponse{
				IsSucceed: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreCheck(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalAPI_GetConfig(t *testing.T) {
	//Initialize()
	tests := []struct {
		name string
		o    *ExternalAPI
		want *webconf.APIServer
	}{
		{
			name: "GetConfig",
			o:    extAPI,
			want: extAPI.apiConf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalAPI.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPI(t *testing.T) {
	//Initialize()
	tests := []struct {
		name string
		want *ExternalAPI
	}{
		{
			name: "ExternalAPI",
			want: &ExternalAPI{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPI(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalAPI_APIs(t *testing.T) {
	//Initialize()

	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		o       *ExternalAPI
		args    args
		wantErr bool
	}{
		{
			name: "ExternalAPI_APIs",
			o:    extAPI,
			args: args{
				c: echoContext,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.GetHealthCheck(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.GetHealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetVersion(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.Version() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := tt.o.CreateTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.CreateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.UpdateTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.UpdateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.DeleteTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.DeleteTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetTemplates(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ExternalAPI.GetTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
