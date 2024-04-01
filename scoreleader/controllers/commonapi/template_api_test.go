package commonapi

import (
	"net/http"
	"testing"

	"stove-gitlab.sginfra.net/backend/template/controllers/context"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
)

func TestCreateTemplate(t *testing.T) {
	Initialize()

	type args struct {
		body string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		/*
			{
				name: "CreateTemplate",
				args: args{
					body: `{"user_id":"unit_test","user_name":"unit_test","role":"admin"}`,
				},
				wantErr: false,
			},
		*/
		{
			name: "1",
			args: args{
				body: ``,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				body: `{}`,
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				body: `{"key":"` + unit_test_id + `","value":"test1"}`,
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				body: `{"key":"` + unit_test_id2 + `","value":"test2"}`,
			},
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				body: `{"key":"` + unit_test_id3 + `","value":"test3"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			echoContext := createEcho(host, http.MethodPost, tt.args.body, headers)
			ctx := base.GetContext(echoContext).(*context.TemplateContext)
			if err := CreateTemplate(ctx); (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTemplate(t *testing.T) {
	Initialize()

	type args struct {
		body string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				body: ``,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				body: `{"key":"unit_test"}`,
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				body: `{"key":"unit_test", "value":"test11"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			echoContext := createEcho(host, http.MethodPut, tt.args.body, headers)
			ctx := base.GetContext(echoContext).(*context.TemplateContext)

			if err := UpdateTemplate(ctx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	Initialize()

	type args struct {
		KeyID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				KeyID: "",
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				KeyID: "unit_test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			echoContext := createEcho(host, http.MethodGet, "", headers)
			ctx := base.GetContext(echoContext).(*context.TemplateContext)
			ctx.SetParam("key", tt.args.KeyID)

			if err := GetTemplate(ctx); (err != nil) != tt.wantErr {
				t.Errorf("GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTemplates(t *testing.T) {
	Initialize()

	type args struct {
		PageIndex string
		PageSize  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				PageIndex: "0",
				PageSize:  "0",
			},
			wantErr: false,
		},
		{
			name: "1",
			args: args{
				PageIndex: "1",
				PageSize:  "100",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			echoContext := createEcho(host, http.MethodGet, "", headers)
			ctx := base.GetContext(echoContext).(*context.TemplateContext)
			ctx.SetParam("page_index", tt.args.PageIndex)
			ctx.SetParam("page_size", tt.args.PageSize)

			if err := GetTemplates(ctx); (err != nil) != tt.wantErr {
				t.Errorf("GetTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTemplate(t *testing.T) {
	Initialize()

	type args struct {
		body string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				body: ``,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				body: `{}`,
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				body: `{"key":"` + unit_test_id + `"}`,
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				body: `{"key":"` + unit_test_id2 + `"}`,
			},
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				body: `{"key":"` + unit_test_id3 + `"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			echoContext := createEcho(host, http.MethodDelete, tt.args.body, headers)
			ctx := base.GetContext(echoContext).(*context.TemplateContext)

			if err := DeleteTemplate(ctx); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
