package resultcode

import (
	"reflect"
	"testing"
)

func TestGetResultCodeMap(t *testing.T) {
	tests := []struct {
		name string
		want map[int]string
	}{
		{
			name: "1",
			want: GetResultCodeMap(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetResultCodeMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResultCodeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
