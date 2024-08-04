package hive

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HandleType interface {
	string | func(*fiber.Ctx) error | map[string]interface{}
}

type Controller struct {
	routes []func(app *GoNest)
}

func (c *Controller) appendRoute(route func(app *GoNest)) {
	c.routes = append(c.routes, route)
}

func CreateController() (controller Controller) {
	controller = Controller{}
	return controller
}

func (c *Controller) Get(
	path string,
	handler func(*fiber.Ctx) any,
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result := handler(ctx)

		if result == nil {
			return nil
		}

		//get type of result
		var text = fmt.Sprintf("%T", result)

		print("Type: ", text)

		switch v := result.(type) {
		case Map:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		case error:
			return ctx.SendString(v.Error())
		}

		return nil
	}

	c.appendRoute(func(app *GoNest) {
		app.fiber.Get(path, handler2)
	})
}

func (c *Controller) Post(
	path string,
	handler func(*fiber.Ctx) interface{},
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result := handler(ctx)

		if result == nil {
			return nil
		}

		switch v := result.(type) {
		case map[string]interface{}:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		case error:
			return ctx.SendString(v.Error())
		}

		return nil
	}

	c.appendRoute(func(app *GoNest) {
		app.fiber.Post(path, handler2)
	})
}

func (c *Controller) Put(
	path string,
	handler func(*fiber.Ctx) interface{},
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result := handler(ctx)

		if result == nil {
			return nil
		}

		switch v := result.(type) {
		case map[string]interface{}:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		case error:
			return ctx.SendString(v.Error())
		}

		return nil
	}

	c.appendRoute(func(app *GoNest) {
		app.fiber.Put(path, handler2)
	})
}

func (c *Controller) Delete(
	path string,
	handler func(*fiber.Ctx) interface{},
) {

	var handler2 = func(ctx *fiber.Ctx) error {
		result := handler(ctx)

		if result == nil {
			return nil
		}

		switch v := result.(type) {
		case map[string]interface{}:
			return ctx.JSON(v)
		case string:
			return ctx.SendString(v)
		case error:
			return ctx.SendString(v.Error())
		}

		return nil
	}

	c.appendRoute(func(app *GoNest) {
		app.fiber.Delete(path, handler2)
	})
}
