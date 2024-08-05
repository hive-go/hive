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

	var handler2 = func(ctx *fiber.Ctx) error {
		result, err := handler(ctx)

		if result == nil {
			return nil
		}

		if err != nil {
			return ctx.SendString(err.Error())
		}

		switch v := result.(type) {
		case Map:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		}

		return nil
	}

	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Get(prefix+path, handler2)
	})
}

func (c *Controller) Post(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result, err := handler(ctx)

		if err != nil {
			//set status code and error
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

	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Post(prefix+path, handler2)
	})
}

func (c *Controller) Put(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result, err := handler(ctx)

		if result == nil {
			return nil
		}

		if err != nil {
			return ctx.SendString(err.Error())
		}

		switch v := result.(type) {
		case map[string]interface{}:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		}

		return nil
	}

	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Put(prefix+path, handler2)
	})
}

func (c *Controller) Delete(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result, err := handler(ctx)

		if result == nil {
			return nil
		}

		if err != nil {
			return ctx.SendString(err.Error())
		}

		switch v := result.(type) {
		case map[string]interface{}:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		}

		return nil
	}

	c.appendRoute(func(app *fiber.App, prefix string) {
		app.Delete(prefix+path, handler2)
	})
}
