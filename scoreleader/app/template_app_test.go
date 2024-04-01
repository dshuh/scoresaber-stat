package app

import (
	"reflect"
	"sync"
	"testing"

	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	webconf "stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

const (
	configFile = "../etc/conf/config.unittest.yml"
)

func TestTemplateApp_Init(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name    string
		o       *TemplateApp
		args    args
		wantErr bool
	}{
		{
			name: "TemplateApp_Init",
			args: args{
				configFile: configFile,
			},
			o: &TemplateApp{
				BaseApp:    base.BaseApp{},
				conf:       &config.TemplateConfig{},
				configFile: configFile,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Init(tt.args.configFile); (err != nil) != tt.wantErr {
				t.Errorf("TemplateApp.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateApp_Run(t *testing.T) {
	type args struct {
		wg *sync.WaitGroup
	}
	tests := []struct {
		name    string
		o       *TemplateApp
		args    args
		wantErr bool
	}{
		{
			name: "TemplateApp_Run",
			args: args{
				wg: &sync.WaitGroup{},
			},
			o: &TemplateApp{
				BaseApp:    base.BaseApp{},
				conf:       &config.TemplateConfig{},
				configFile: configFile,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Run(tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("TemplateApp.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateApp_GetConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		o          *TemplateApp
		want       *webconf.Config
	}{
		{
			name:       "TemplateApp_GetConfig",
			configFile: configFile,
			o: &TemplateApp{
				BaseApp:    base.BaseApp{},
				conf:       &config.TemplateConfig{},
				configFile: configFile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.o.Init(tt.configFile); err != nil {
				t.Errorf("TemplateApp.Init() error = %v", err)
			}

			tt.want = &tt.o.conf.Config

			if got := tt.o.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TemplateApp.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateApp_CreateTable(t *testing.T) {
	tests := []struct {
		name    string
		o       *TemplateApp
		wantErr bool
	}{
		{
			name: "TemplateApp_CreateTable",
			o: &TemplateApp{
				BaseApp:    base.BaseApp{},
				conf:       &config.TemplateConfig{},
				configFile: configFile,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.CreateTable(); (err != nil) != tt.wantErr {
				t.Errorf("TemplateApp.CreateTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateApp_InitTable(t *testing.T) {
	tests := []struct {
		name    string
		o       *TemplateApp
		wantErr bool
	}{
		{
			name: "TemplateApp_CreateTable",
			o: &TemplateApp{
				BaseApp:    base.BaseApp{},
				conf:       &config.TemplateConfig{},
				configFile: configFile,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.InitTable(); (err != nil) != tt.wantErr {
				t.Errorf("TemplateApp.InitTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateApp_Execute(t *testing.T) {
	o := &TemplateApp{
		BaseApp:    base.BaseApp{},
		conf:       &config.TemplateConfig{},
		configFile: configFile,
	}

	tests := []struct {
		name    string
		o       *TemplateApp
		wantErr bool
	}{
		{
			name:    "createtable",
			o:       o,
			wantErr: false,
		},
		{
			name:    "test",
			o:       o,
			wantErr: false,
		},
		{
			name:    "inittable",
			o:       o,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.o.Init(o.configFile); err != nil {
				t.Errorf("TemplateApp.Init() error = %v", err)
			}

			if tt.name == "createtable" {
				*command = "createtable"
			} else if tt.name == "test" {
				*command = "test"
			} else if tt.name == "inittable" {
				*command = "inittable"
				*serviceType = "test"
				*serviceID = "test"
			}

			if err := tt.o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("TemplateApp.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *TemplateApp
		wantErr bool
	}{
		{
			name: "TemplateApp_NewApp",
			args: args{
				version: "",
			},
			want:    &TemplateApp{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApp(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want = got
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
