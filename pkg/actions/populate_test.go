package actions

import (
	"reflect"
	"testing"

	"github.com/stakater/Chowkidar/pkg/config"
)

func TestPopulateFromConfig(t *testing.T) {
	type args struct {
		configActions []config.Action
		criterion     config.Criterion
	}
	tests := []struct {
		name string
		args args
		want []Action
	}{
		{
			name: "PopulateSlackAction",
			args: args{
				configActions: []config.Action{
					config.Action{
						Name: "slack",
						Params: map[interface{}]interface{}{
							"token":   "123",
							"channel": "channelName",
						},
					},
				},
				criterion: config.Criterion{},
			},
			want: []Action{
				MapToAction("slack"),
			},
		},
		{
			name: "PopulateSlackActionError",
			args: args{
				configActions: []config.Action{
					config.Action{
						Name: "slack",
						Params: map[interface{}]interface{}{
							"tok":  "123",
							"chan": "channelName",
						},
					},
				},
				criterion: config.Criterion{},
			},
			want: []Action{
				MapToAction("slack"),
			},
		},
		{
			name: "PopulateDefaultAction",
			args: args{
				configActions: []config.Action{
					config.Action{
						Name: "default",
					},
				},
				criterion: config.Criterion{},
			},
			want: []Action{
				MapToAction("default"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PopulateFromConfig(tt.args.configActions, tt.args.criterion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PopulateFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
