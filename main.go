package hive

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Version of current hive package
const Version = "1.0.1"

type GoNest struct {
	fiberApp *fiber.App
	modules  []Module
}

type Config = fiber.Config

func New(config ...Config) (instance *GoNest) {
	fiber := fiber.New(config...)
	app := GoNest{}
	app.fiberApp = fiber
	return &app
}

func (n *GoNest) Listen(port string) {

	var Green = "\033[32m"
	var Reset = "\033[0m"

	now := time.Now()

	exec.Command("clear")

	fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Starting Hive application..." + Reset)

	for _, module := range n.modules {

		var prefixModule = ""

		if module.config.Prefix != "" {
			prefixModule = module.config.Prefix
		}

		for _, controller := range module.controllers {

			var prefixController = ""

			if controller.config.Prefix != "" {
				prefixController = controller.config.Prefix
			}
			for _, generateRoute := range controller.routes {

				generateRoute(n.fiberApp, prefixModule+prefixController)
			}
		}
	}

	fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Application started on port " + port + Reset)
	n.fiberApp.Listen("127.0.0.1:" + port)
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
