package commonapi

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/controllers/resultcode"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
)

const (
	configFile = "../../etc/conf/config.unittest.yml"

	unit_test_id  = "unit_test"
	unit_test_id2 = "unit_test2"
	unit_test_id3 = "unit_test3"
)

var headers = []string{"content-type", "application/json"}

var host = "http://localhost:10473/"

var conf *config.TemplateConfig
var templateCtx *context.TemplateContext

var TestServiceID string

type args struct {
	ctx *context.TemplateContext
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
				ctx: templateCtx,
			},
			wantErr: false,
		})
	}
	return tests
}

func getUploadTests(caseCount int, fileParam string, filePath []string) []test {
	var tests []test
	for i := 1; i <= caseCount; i++ {
		if i == 1 {
			UploadInitialize(fileParam+"sfs", filePath[i-1])
			tests = append(tests, test{
				name: strconv.Itoa(i),
				args: args{
					ctx: templateCtx,
				},
				wantErr: false,
			})
		} else {
			UploadInitialize(fileParam, filePath[i-1])
			tests = append(tests, test{
				name: strconv.Itoa(i),
				args: args{
					ctx: templateCtx,
				},
				wantErr: false,
			})
		}
	}
	return tests
}

func TestMain(t *testing.T) {
	Initialize()
}

func Initialize() error {
	conf = config.GetInstance(configFile)
	InitContext()
	TestServiceID = `unittest_` + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return nil
}
func InitContext() error {
	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	context.AppendRequestParameter()
	echoContext := createEcho(host, http.MethodPost, "", headers)
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
	echoContext := e.NewContext(req, rec)

	if err := base.SetContext(echoContext, &conf.Config, context.NewTemplateContext); err != nil {
		log.Error(err)
	}

	return echoContext
}

func UploadInitialize(fileParam, filePath string) error {
	conf = config.GetInstance(configFile)
	base.AppendResultCodeText(&resultcode.TemplateResultCodeMap)
	context.AppendRequestParameter()
	echoContext := createUploadEcho(host, http.MethodPost, fileParam, filePath)
	templateCtx = base.GetContext(echoContext).(*context.TemplateContext)
	TestServiceID = `unittest_` + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return nil
}

func createUploadEcho(url, method, fileParam, filePath string) echo.Context {

	e := echo.New()

	file, err := os.Open(filePath)
	defer file.Close()

	buf := new(bytes.Buffer)
	mr := multipart.NewWriter(buf)
	part, err := mr.CreateFormFile(fileParam, filepath.Base(filePath))
	if err != nil {
		return nil
	}
	_, err = io.Copy(part, file)
	err = mr.Close()

	req := httptest.NewRequest(method, url, buf)
	req.Header.Set("Content-Type", mr.FormDataContentType())

	rec := httptest.NewRecorder()
	echoContext := e.NewContext(req, rec)

	if err := base.SetContext(echoContext, &conf.Config, context.NewTemplateContext); err != nil {
		log.Error(err)
	}

	return echoContext
}
func TestGetHealthCheck(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetHealthCheck",
			args: args{
				c: createEcho(host, http.MethodGet, "", headers),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetHealthCheck(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetHealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	type args struct {
		c       echo.Context
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetVersion",
			args: args{
				c:       createEcho(host, http.MethodGet, "", headers),
				version: "m1.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetVersion(tt.args.c, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	type args struct {
		c       echo.Context
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetStatus",
			args: args{
				c: createEcho(host, http.MethodGet, "", headers),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetStatus(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRealIP(t *testing.T) {
	type args struct {
		c       echo.Context
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetRealIP",
			args: args{
				c:       createEcho(host, http.MethodGet, "", headers),
				version: "m1.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetRealIP(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetRealIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
