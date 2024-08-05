package hive

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GoNest struct {
	fiber   *fiber.App
	modules []Module
	Ctx     *fiber.Ctx
}

func New() (instance *GoNest) {
	fiber := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app := GoNest{}
	app.fiber = fiber
	return &app
}

func (n *GoNest) Listen(port string) {

	var Green = "\033[32m"
	var Reset = "\033[0m"

	now := time.Now()

	exec.Command("clear")

	fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Starting Hive application..." + Reset)

	for _, module := range n.modules {
		for _, controller := range module.controllers {
			for _, generateRoute := range controller.routes {
				generateRoute(n)
			}
		}
	}

	fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Application started on port " + port + Reset)
	n.fiber.Listen("127.0.0.1:" + port)
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
