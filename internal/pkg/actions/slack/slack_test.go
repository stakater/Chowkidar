package slack

import (
	"log"
	"testing"

	"github.com/stakater/Chowkidar/internal/pkg/config"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	configFilePath   = "../../../../configs/testConfigs/CorrectSlackConfig.yaml"
	configuration, _ = config.ReadConfig(configFilePath)
)

type SlackMock struct {
}

func (s *SlackMock) SendNotification(message string) error {
	log.Print(message)
	return nil
}
func TestSlack_Init(t *testing.T) {
	type fields struct {
		Token     string
		Channel   string
		Criterion config.Criterion
	}
	type args struct {
		params    map[interface{}]interface{}
		criterion config.Criterion
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "MissingSlackToken",
			args: args{
				params: map[interface{}]interface{}{
					"token":   "",
					"channel": "channelName",
				},
				criterion: config.Criterion{},
			},
			wantErr: true,
		},
		{
			name: "CorrectScenario",
			args: args{
				params: map[interface{}]interface{}{
					"token":   "123",
					"channel": "channelName",
				},
				criterion: config.Criterion{},
			},
			fields: fields{
				Token:     "123",
				Channel:   "channelName",
				Criterion: config.Criterion{},
			},
		},

		{
			name: "ErrorInDecoding",
			args: args{
				params: map[interface{}]interface{}{
					"tokens":  "123",
					"channel": "channelName",
				},
				criterion: config.Criterion{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Slack{
				Token:     tt.fields.Token,
				Channel:   tt.fields.Channel,
				Criterion: tt.fields.Criterion,
			}
			if err := s.Init(tt.args.params, tt.args.criterion); (err != nil) != tt.wantErr {
				t.Errorf("Slack.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlack_ObjectCreated(t *testing.T) {
	type fields struct {
		Token     string
		Channel   string
		Criterion config.Criterion
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ObjectCreated",
			args: args{
				obj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test",
						Namespace: "asd",
					},
				},
			},
			fields: fields{
				Token:   configuration.Controllers[0].Actions[0].Params["token"].(string),
				Channel: configuration.Controllers[0].Actions[0].Params["channel"].(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackMock{}
			s.SendNotification("Sending create notification")
		})
	}
}

func TestSlack_ObjectUpdated(t *testing.T) {
	type fields struct {
		Token     string
		Channel   string
		Criterion config.Criterion
	}
	type args struct {
		oldObj interface{}
		newObj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ObjectUpdated",
			args: args{
				oldObj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test-old",
						Namespace: "asd",
					},
				},
				newObj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test-new",
						Namespace: "asd",
					},
				},
			},
			fields: fields{
				Token:   configuration.Controllers[0].Actions[0].Params["token"].(string),
				Channel: configuration.Controllers[0].Actions[0].Params["channel"].(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackMock{}
			s.SendNotification("Sending update notification")
		})
	}
}

func TestSlack_ObjectDeleted(t *testing.T) {
	type fields struct {
		Token     string
		Channel   string
		Criterion config.Criterion
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ObjectDeleted",
			args: args{
				obj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test",
						Namespace: "asd",
					},
				},
			},
			fields: fields{
				Token:   configuration.Controllers[0].Actions[0].Params["token"].(string),
				Channel: configuration.Controllers[0].Actions[0].Params["channel"].(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackMock{}
			s.SendNotification("Sending delete notification")
		})
	}
}
