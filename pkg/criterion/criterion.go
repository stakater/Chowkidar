package criterion

import (
	"reflect"

	"github.com/stakater/Chowkidar/pkg/config"
)

func MatchFuncSingle(mapFunc interface{}, criterion config.Criterion) func(interface{}) {
	return func(value interface{}) {
		//TODO: Create a criterion matcher and use it here
		if matchesCriterion(value, criterion) {
			funcValue := reflect.ValueOf(mapFunc)

			funcValue.Call([]reflect.Value{
				reflect.ValueOf(value),
			})
		}
	}
}

// TODO: Create a generic arbitrary func for any number of arguments or specifically 1 and 2 arguments
func MatchFuncMulti(mapFunc interface{}, criterion config.Criterion) func(interface{}, interface{}) {
	return func(value1 interface{}, value2 interface{}) {
		//TODO: Create a criterion matcher and use it here
		if matchesCriterion(value1, criterion) {
			funcValue := reflect.ValueOf(mapFunc)

			funcValue.Call([]reflect.Value{
				reflect.ValueOf(value1),
				reflect.ValueOf(value2),
			})
		}
	}
}
