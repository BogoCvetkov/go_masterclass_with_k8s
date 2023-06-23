package controller

import (
	"github.com/BogoCvetkov/go_mastercalss/db"
)

type ControllerMap struct {
	Account  *accController
	Transfer *transferController
	User     *userController
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
		Transfer: &transferController{
			apiController: baseCtr,
		},
		User: &userController{
			apiController: baseCtr,
		},
	}

	return controllerMap
}
