package hive

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Version of current hive package
const Version = "1.0.1"

type GoNest struct {
	*fiber.App
	modules []Module
	config  Config
}

type Config struct {
	FiberConfig   fiber.Config
	SwaggerConfig SwaggerConfig
}

var Green = "\033[32m"
var Reset = "\033[0m"

func New(config Config) (instance *GoNest) {
	now := time.Now()
	fmt.Println(Green + "\n[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Starting Hive application..." + Reset)

	fiber := fiber.New(config.FiberConfig)

	app := GoNest{}
	app.App = fiber
	app.config = config
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

	now := time.Now()
	fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive] Application started on Address " + addr + Reset)

	if n.config.SwaggerConfig.Enabled {
		// generateSwagger(n)
		GenerateSwaggerV2(n)
		path := "/api"

		if n.config.SwaggerConfig.Path != "" {
			path = n.config.SwaggerConfig.Path
		}

		fullPath := path + "/*"

		n.App.Get(fullPath, swagger.New(swagger.Config{
			ConfigURL: "/swagger",
			URL:       "/swagger",
		}))

		n.App.Get("/swagger", func(c *fiber.Ctx) error {
			return c.SendFile("./swagger.json")
		})
	}

	n.modules = nil
	n.config = Config{}

	n.App.Listen(addr)
}

func (n *GoNest) AddModule(module Module) {
	n.modules = append(n.modules, module)
}

func (n *GoNest) CreateModule() (module Module) {
	module = Module{}
	return module
}
