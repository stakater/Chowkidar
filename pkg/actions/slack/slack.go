package slack

import (
	"fmt"
	"log"
	"strings"

	"github.com/asiyani/slack"
	"github.com/mitchellh/mapstructure"
	"github.com/stakater/Chowkidar/pkg/config"
)

type Slack struct {
	Token     string
	Channel   string
	Criterion config.Criterion
}

func (s *Slack) Init(params map[interface{}]interface{}, criteron config.Criterion) error {
	s.Criterion = criteron
	err := mapstructure.Decode(params, &s)
	if err != nil {
		panic(err)
	}
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf("Missing slack token or channel")
	}
	return nil
}

func (s *Slack) ObjectCreated(obj interface{}) {
	sendSlackNotification(s, "Created")

}

func (s *Slack) ObjectDeleted(obj interface{}) {
	sendSlackNotification(s, "Deleted")
}

func (s *Slack) ObjectUpdated(oldObj, newObj interface{}) {
	sendSlackNotification(s, "Updated")
}
func sendSlackNotification(s *Slack, input string) {

	// fmt.Print(s.Token, s.Channel)

	api := slack.New(s.Token)
	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{prepareMessage(s)}
	params.AsUser = false

	channelID, timestamp, err := api.PostMessage(s.Channel, "A controller has been "+input, params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
func prepareMessage(s *Slack) slack.Attachment {
	return slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Chowkidar",
				Value: strings.Join(s.Criterion.Identifiers, " "+s.Criterion.Operator+" "),
			},
		},
	}
}
