package actions

import (
	"log"

	"github.com/stakater/Chowkidar/internal/pkg/actions/slack"
	"github.com/stakater/Chowkidar/internal/pkg/config"
	"k8s.io/api/core/v1"
)

func assertActionImplementations() {
	var _ Action = (*Default)(nil)
	var _ Action = (*slack.Slack)(nil)
}

const (
	DefaultAction = "default"
)

// Action interface so that other actions like slack can implement this
type Action interface {
	Init(map[interface{}]interface{}, config.Criterion) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}

// MapToAction maps the action name to the actual action type
func MapToAction(actionName string) Action {
	action, ok := Map[actionName]
	if !ok {
		return Map[DefaultAction]
	}
	return action
}

var Map = map[string]Action{
	"default": &Default{},
	"slack":   &slack.Slack{},
}

// Default class with empty implementations for any action that we dont support currently
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init(params map[interface{}]interface{}, criterion config.Criterion) error {
	return nil
}

// ObjectCreated Do nothing for default handler
func (d *Default) ObjectCreated(obj interface{}) {
	message := "Resource block not found for Pod: `" + obj.(*v1.Pod).Name + "` in Namespace: `" + obj.(*v1.Pod).Namespace + "`"
	log.Printf(message)
}

// ObjectDeleted Do nothing for default handler
func (d *Default) ObjectDeleted(obj interface{}) {

}

// ObjectUpdated Do nothing for default handler
func (d *Default) ObjectUpdated(oldObj, newObj interface{}) {
	message := "Resource block not found for Pod: `" + oldObj.(*v1.Pod).Name + "` in Namespace: `" + oldObj.(*v1.Pod).Namespace + "`"
	log.Printf(message)
}
