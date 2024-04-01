package config

import (
	"reflect"
	"testing"
)

const (
	configFile = "../etc/conf/config.unittest.yml"
)

func TestGetInstance(t *testing.T) {
	type args struct {
		filepath []string
	}
	tests := []struct {
		name string
		args args
		want *TemplateConfig
	}{
		{
			name: "1",
			args: args{
				filepath: []string{"invalid_" + configFile},
			},
		},
		{
			name: "2",
			args: args{
				filepath: []string{configFile},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstance(tt.args.filepath...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}
