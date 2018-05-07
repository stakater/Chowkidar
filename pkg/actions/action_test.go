package actions

import (
	"reflect"
	"testing"

	"github.com/stakater/Chowkidar/pkg/actions/slack"
)

func TestMapToAction(t *testing.T) {
	type args struct {
		actionName string
	}
	tests := []struct {
		name string
		args args
		want Action
	}{
		{
			name: "MapDefault",
			args: args{
				actionName: "default",
			},
			want: &Default{},
		},
		{
			name: "MapSlack",
			args: args{
				actionName: "slack",
			},
			want: &slack.Slack{},
		},
		{
			name: "MapDefaultFromError",
			args: args{
				actionName: "",
			},
			want: &Default{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToAction(tt.args.actionName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
