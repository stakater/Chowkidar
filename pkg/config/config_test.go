package config

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "TestingWithCorrectValues",
			args: args{filePath: "../../configs/testConfigs/CorrectConfig.yaml"},
			want: Config{
				Controllers: []Controller{
					Controller{
						Type: "pods",
						WatchCriterion: Criterion{
							Operator: "and",
							Identifiers: []string{
								"resourceExists",
								"healthCheckExists",
							},
						},
						Actions: []Action{
							Action{
								Name: "slack",
								Params: map[interface{}]interface{}{
									"token":   "123",
									"channel": "channelName",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TestingWithEmptyFile",
			args: args{filePath: "../../configs/testConfigs/Empty.yaml"},
			want: Config{},
		},
		{
			name: "TestingWithNoAction",
			args: args{filePath: "../../configs/testConfigs/NoActions.yaml"},
			want: Config{
				Controllers: []Controller{
					Controller{
						Type: "pods",
						WatchCriterion: Criterion{
							Operator: "and",
							Identifiers: []string{
								"resourceExists",
								"healthCheckExists",
							},
						},
					},
				},
			},
		},
		{
			name:    "TestingWithFileNotPresent",
			args:    args{filePath: "../../configs/testConfigs/FileNotFound.yaml"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
