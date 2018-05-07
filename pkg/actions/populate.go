package actions

import (
	"log"

	"github.com/stakater/Chowkidar/pkg/config"
)

// PopulateFromConfig populates the actions for a specific controller from config
func PopulateFromConfig(configActions []config.Action, criterion config.Criterion) []Action {
	var populatedActions []Action
	for _, configAction := range configActions {
		actionToAdd := MapToAction(configAction.Name)
		err := actionToAdd.Init(configAction.Params, criterion)
		if err != nil {
			log.Println(err)
		}
		populatedActions = append(populatedActions, actionToAdd)
	}
	return populatedActions
}
