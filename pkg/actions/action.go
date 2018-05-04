package actions

import "github.com/stakater/Chowkidar/pkg/config"

// Action interface so that other actions like slack can implement this
type Action interface {
	Init(map[interface{}]interface{}, config.Criterion) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
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

}

// ObjectDeleted Do nothing for default handler
func (d *Default) ObjectDeleted(obj interface{}) {

}

// ObjectUpdated Do nothing for default handler
func (d *Default) ObjectUpdated(oldObj, newObj interface{}) {

}
