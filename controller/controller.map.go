package controller

import (
	"github.com/BogoCvetkov/go_mastercalss/db"
)

type ControllerMap struct {
	Account *accController
}

// Initialize Controllers
func InitControllers(s *db.Store) *ControllerMap {
	var controllerMap *ControllerMap

	baseCtr := &apiController{
		store: s,
	}

	controllerMap = &ControllerMap{
		Account: &accController{
			apiController: baseCtr,
		},
	}

	return controllerMap
}
