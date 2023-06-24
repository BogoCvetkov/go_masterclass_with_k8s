package api

import (
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
)

type ControllerMap struct {
	Account  *accController
	Transfer *transferController
	User     *userController
}

// Initialize Controllers
func InitControllers(s interfaces.IServer) *ControllerMap {
	var controllerMap *ControllerMap

	baseCtr := &apiController{
		server: s,
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
