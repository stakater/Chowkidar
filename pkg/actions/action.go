package actions

import "github.com/stakater/Chowkidar/pkg/config"

type Action interface {
	Init(map[interface{}]interface{}, config.Criterion) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init(params map[interface{}]interface{}, criterion config.Criterion) error {
	return nil
}

func (d *Default) ObjectCreated(obj interface{}) {

}

func (d *Default) ObjectDeleted(obj interface{}) {

}

func (d *Default) ObjectUpdated(oldObj, newObj interface{}) {

}
