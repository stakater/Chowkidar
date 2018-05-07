package slack

import (
	"errors"
	"log"

	"github.com/asiyani/slack"
	"github.com/mitchellh/mapstructure"
	"github.com/stakater/Chowkidar/pkg/config"
	"k8s.io/api/core/v1"
)

// Slack action class implementing the Action interface
type Slack struct {
	Token     string
	Channel   string
	Criterion config.Criterion
}

// Init initializes the Slack Configuration like token and channel
func (s *Slack) Init(params map[interface{}]interface{}, criterion config.Criterion) error {
	s.Criterion = criterion
	err := mapstructure.Decode(params, &s) //Converts the params to slack struct fields
	if err != nil {
		return err
	}
	if s.Token == "" || s.Channel == "" {
		return errors.New("Missing slack token or channel")
	}
	return nil
}

// ObjectCreated sending SlackNotification when an object is created
func (s *Slack) ObjectCreated(obj interface{}) {
	message := "Resource block not found for Pod: `" + obj.(*v1.Pod).Name + "` in Namespace: `" + obj.(*v1.Pod).Namespace + "`"
	sendSlackNotification(s, message)
}

// ObjectDeleted sending SlackNotification when an object is deleted
func (s *Slack) ObjectDeleted(obj interface{}) {

}

// ObjectUpdated sending SlackNotification when an object is updated
func (s *Slack) ObjectUpdated(oldObj, newObj interface{}) {
	message := "Resource block not found for Pod: `" + oldObj.(*v1.Pod).Name + "`"
	sendSlackNotification(s, message)
}

// sends the Notification based on the event
func sendSlackNotification(s *Slack, message string) {
	api := slack.New(s.Token)
	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{prepareMessage(s, message)}
	params.AsUser = false

	_, _, err := api.PostMessage(s.Channel, "Chowkidar Alert", params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to Slack Channel `%s`", s.Channel)
}

// Prepares the attachments to send in POST request
func prepareMessage(s *Slack, message string) slack.Attachment {
	return slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: message,
			},
		},
	}
}
