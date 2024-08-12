package hive

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Version of current hive package
const Version = "1.0.1"

type GoNest struct {
	*fiber.App
	modules []Module
}

type Config = fiber.Config

var Green = "\033[32m"
var Reset = "\033[0m"

func New(config ...Config) (instance *GoNest) {
	now := time.Now()
	fmt.Println(Green + "\n[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Starting Hive application..." + Reset)
	fiber := fiber.New(config...)
	app := GoNest{}
	app.App = fiber
	return &app
}

func (n *GoNest) Listen(addr string) {

	exec.Command("clear")

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

			for _, route := range controller.routes {
				route.function(n.App, prefixModule+prefixController)
			}
		}

	}

	generateSwagger(&n.modules)
	now := time.Now()
	fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Application started on port " + addr + Reset)

	n.App.Get("/api/*", swagger.New(swagger.Config{
		ConfigURL: "/swagger",
		URL:       "/swagger",
	}))

	//give file at root ./swagger.json
	n.App.Get("/swagger", func(c *fiber.Ctx) error {
		return c.SendFile("./swagger.json")
	})

	//open browser in swagger
	open(addr + "/api")

	n.App.Listen(addr)
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
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
