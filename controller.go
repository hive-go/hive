package hive

import (
	"github.com/gofiber/fiber/v2"
)

type HandleType interface {
	string | func(*fiber.Ctx) error | map[string]interface{}
}

type Controller struct {
	routes []func(app *fiber.App, prefix string)
	config ControllerConfig
}

type ControllerConfig struct {
	Prefix string
}

func (c *Controller) SetConfig(config ControllerConfig) {
	c.config = config
}

func (c *Controller) appendRoute(route func(app *fiber.App, prefix string)) {
	c.routes = append(c.routes, route)
}

func CreateController() (controller Controller) {
	controller = Controller{}
	return controller
}

func (c *Controller) Get(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Get(prefix+path, generateHandlerOfCallback(handler))
	})
}

func (c *Controller) Post(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Post(prefix+path, generateHandlerOfCallback(handler))
	})
}

func (c *Controller) Put(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Put(prefix+path, generateHandlerOfCallback(handler))
	})
}

func (c *Controller) Delete(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Delete(prefix+path, generateHandlerOfCallback(handler))
	})
}

func generateHandlerOfCallback(callback func(ctx *fiber.Ctx) (any, error)) func(ctx *fiber.Ctx) error {

	var handler = func(ctx *fiber.Ctx) error {
		result, err := callback(ctx)

		if err != nil {
			return err
		}

		switch v := result.(type) {
		case map[string]interface{}:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		}

		return ctx.JSON(result)
	}

	return handler
}
