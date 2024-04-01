package internalapi

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

//var version = "m1.0"
//var prefixUrl = host + version

// var healthcheckUrl = host + "/healthcheck/_check"
// var versionUrl = prefixUrl + "/version"

var intAPI *InternalAPI
var conf *config.TemplateConfig
var echoContext echo.Context
var templateCtx *context.TemplateContext

func TestMain(t *testing.T) {
	Initialize()
}

func Initialize() error {
	conf = config.GetInstance(configFile)

	intAPI = NewAPI()
	intAPI.apiConf = intAPI.GetConfig()
	intAPI.echo = echo.New()

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

func TestInternalAPI_Init(t *testing.T) {
	//Initialize()
	if intAPI.apiConf.Routes[0:1] != "/" {
		intAPI.apiConf.Routes = "../../etc/conf/" + intAPI.apiConf.Routes
	}

	type args struct {
		e *echo.Echo
	}
	tests := []struct {
		name    string
		o       *InternalAPI
		args    args
		wantErr bool
	}{
		{
			name: "InternalAPI_Init",
			o:    intAPI,
			args: args{
				e: intAPI.echo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Init(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreCheck(t *testing.T) {
	//Initialize()
	type args struct {
		c    echo.Context
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				c:    echoContext,
				path: "login",
			},
		},
		{
			name: "2",
			args: args{
				c:    echoContext,
				path: "",
			},
		},
		{
			name: "3",
			args: args{
				c:    nil,
				path: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.c != nil {
				tt.args.c.SetPath(tt.args.path)
			}
			if got := PreCheck(tt.args.c); (got.IsSucceed != true) != tt.wantErr {
				t.Errorf("PreCheck() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestInternalAPI_GetConfig(t *testing.T) {
	//Initialize()
	tests := []struct {
		name string
		o    *InternalAPI
		want *webconf.APIServer
	}{
		{
			name: "GetConfig",
			o:    intAPI,
			want: intAPI.apiConf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InternalAPI.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPI(t *testing.T) {
	//Initialize()
	tests := []struct {
		name string
		want *InternalAPI
	}{
		{
			name: "InternalAPI",
			want: &InternalAPI{},
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
func TestInternalAPI_APIs(t *testing.T) {
	//Initialize()

	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		o       *InternalAPI
		args    args
		wantErr bool
	}{
		{
			name: "InternalAPI_APIs",
			o:    intAPI,
			args: args{
				c: echoContext,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// System
			if err := tt.o.GetHealthCheck(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.GetHealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetVersion(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.Version() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetRealIP(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.Version() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := tt.o.CreateTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.CreateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.UpdateTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.UpdateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.DeleteTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.DeleteTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.o.GetTemplates(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalAPI.GetTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
