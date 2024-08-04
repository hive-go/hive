package hive

import (
	"github.com/gofiber/fiber/v2"
)

type GoNest struct {
	fiber   *fiber.App
	modules []Module
	Ctx     *fiber.Ctx
}

func New() (instance *GoNest) {
	fiber := fiber.New()
	app := GoNest{}
	app.fiber = fiber
	return &app
}

func (n *GoNest) Listen(port string) {

	for _, module := range n.modules {
		for _, controller := range module.controllers {
			for _, generateRoute := range controller.routes {
				generateRoute(n)
			}
		}
	}

	n.fiber.Listen(port)
}

func CreateModule() (module Module) {
	module = Module{}
	return module
}

func CreateService() (service Service) {
	service = Service{}
	return service
}

func (n *GoNest) AddModule(module Module) {
	n.modules = append(n.modules, module)
}

func (n *GoNest) CreateModule() (module Module) {
	module = Module{}
	return module
}
