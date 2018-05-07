package kube

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestMapToRuntimeObject(t *testing.T) {
	type args struct {
		resourceType string
	}
	tests := []struct {
		name string
		args args
		want runtime.Object
	}{
		{
			name: "MapPods",
			args: args{
				resourceType: "pods",
			},
			want: &v1.Pod{},
		},
		{
			name: "MapDefault",
			args: args{
				resourceType: "default",
			},
			want: nil,
		},
		{
			name: "MapDefaultThroughError",
			args: args{
				resourceType: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToRuntimeObject(tt.args.resourceType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToRuntimeObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
