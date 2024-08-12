package hive

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct {
	preRoutes map[string][]func(app *fiber.Ctx) error
	routes    map[string]struct {
		function func(app *fiber.App, prefix string) error
		metadata map[string]interface{}
	}
	metadata   map[string]interface{}
	orderCount int
	orderId    string
	config     ControllerConfig
}

type ControllerConfig struct {
	Prefix string
	Tag    string
}

func (c *Controller) SetConfig(config ControllerConfig) {
	c.config = config
	if config.Tag != "" {
		c.metadata["tag"] = config.Tag
	}
}

func (c *Controller) appendPreRoute(id string, route func(app *fiber.Ctx) error) {
	c.preRoutes[id] = append(c.preRoutes[id], route)
}

func (c *Controller) appendRoute(id string, route func(app *fiber.App, prefix string) error) {
	tempRoute := c.routes[id]
	tempRoute.function = route
	c.routes[id] = tempRoute
}

func CreateController() (controller Controller) {
	controller = Controller{}
	controller.orderCount = 0
	controller.orderId = ""
	controller.config = ControllerConfig{}
	controller.metadata = make(map[string]interface{})
	controller.preRoutes = make(map[string][]func(app *fiber.Ctx) error)
	controller.routes = make(map[string]struct {
		function func(app *fiber.App, prefix string) error
		metadata map[string]interface{}
	})
	return controller
}

func (c *Controller) Use(
	handler func(*fiber.Ctx) (interface{}, error),
) *Controller {

	orderCount := c.orderCount

	if orderCount == 0 {
		c.orderId = uuid.New().String()
	}

	c.appendPreRoute(c.orderId, func(ctx *fiber.Ctx) error {
		_, err := handler(ctx)

		if err != nil {
			return err
		}

		return nil
	})

	c.orderCount++

	return c
}

func (c *Controller) EnableBearerAuth() *Controller {
	orderCount := c.orderCount

	if orderCount == 0 {
		c.orderId = uuid.New().String()
	}

	tempRoute := c.routes[c.orderId]

	if tempRoute.metadata == nil {
		tempRoute.metadata = make(map[string]interface{})
	}

	tempRoute.metadata["bearer"] = true
	c.routes[c.orderId] = tempRoute

	c.orderCount++

	return c
}

var myValidator = XValidator{validator: validate}

func (c *Controller) ParseBody(
	bodyInterface interface{},
) *Controller {
	orderCount := c.orderCount

	if orderCount == 0 {
		c.orderId = uuid.New().String()
	}

	tempRoute := c.routes[c.orderId]

	if tempRoute.metadata == nil {
		tempRoute.metadata = make(map[string]interface{})
	}

	tempRoute.metadata["body"] = bodyInterface

	c.routes[c.orderId] = tempRoute

	//add validator to prehooks

	c.appendPreRoute(c.orderId, func(ctx *fiber.Ctx) error {

		body := reflect.New(reflect.TypeOf(bodyInterface)).Interface()

		if err := ctx.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "BODY_PARSING_ERROR")
		}

		if errs := myValidator.Validate(body); len(errs) > 0 && errs[0].Error {
			errMsgs := make([]string, 0)

			for _, err := range errs {
				errMsgs = append(errMsgs, fmt.Sprintf(
					"[%s]: '%v' | Needs to implement '%s'",
					err.FailedField,
					err.Value,
					err.Tag,
				))
			}

			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: strings.Join(errMsgs, " and "),
			}
		}

		ctx.Locals("body", body)

		return nil
	})

	c.orderCount++

	return c
}

func (c *Controller) Get(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(c.orderId, c.generateFinalHandler(path, handler, "GET"))
}

func (c *Controller) Post(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(c.orderId, c.generateFinalHandler(path, handler, "POST"))
}

func (c *Controller) Put(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(c.orderId, c.generateFinalHandler(path, handler, "PUT"))
}

func (c *Controller) Patch(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(c.orderId, c.generateFinalHandler(path, handler, "PATCH"))
}

func (c *Controller) Delete(
	path string,
	handler func(*fiber.Ctx) (interface{}, error),
) {
	c.appendRoute(c.orderId, c.generateFinalHandler(path, handler, "DELETE"))
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

func (c *Controller) generateFinalHandler(path string, handler func(*fiber.Ctx) (interface{}, error), method string) func(*fiber.App, string) error {

	if c.orderCount == 0 {
		c.orderId = uuid.New().String()
	}

	c.orderCount = 0

	allPrefixHandlers := c.preRoutes[c.orderId]

	orderId := c.orderId

	return func(app *fiber.App, prefix string) error {

		var methodToApply func(path string, handlers ...func(*fiber.Ctx) error) fiber.Router

		fullpath := prefix + path

		switch method {
		case "GET":
			methodToApply = app.Get
		case "POST":
			methodToApply = app.Post
		case "PUT":
			methodToApply = app.Put
		case "PATCH":
			methodToApply = app.Patch
		case "DELETE":
			methodToApply = app.Delete
		}

		now := time.Now()

		methodStringWithSpace := method

		if methodStringWithSpace == "DELETE" {
			methodStringWithSpace = "DEL "
		}

		if len(method) == 3 {
			methodStringWithSpace = method + " "
		}

		var Green = "\033[32m"
		var Reset = "\033[0m"
		fmt.Println(Green + "[Hive] - " + Reset + now.Format("02/01/2006, 15:04:05") + Green + " LOG [Hive]  - " + "Registering " + methodStringWithSpace + " route: " + Reset + fullpath)

		tempRoute := c.routes[orderId]

		if tempRoute.metadata == nil {
			tempRoute.metadata = make(map[string]interface{})
		}

		tempRoute.metadata["path"] = fullpath
		tempRoute.metadata["method"] = method
		tempRoute.metadata["operationId"] = orderId
		tempRoute.metadata["parameters"] = getParameters(fullpath)

		c.routes[orderId] = tempRoute

		methodToApply(prefix+path, func(ctx *fiber.Ctx) error {
			for _, handler := range allPrefixHandlers {
				err := handler(ctx)
				if err != nil {
					return err
				}
			}

			return generateHandlerOfCallback(handler)(ctx)
		})

		return nil
	}
}

func getParameters(path string) []string {
	parameters := []string{}

	// /domain/:id/validate

	currentParameter := ""
	startGettingParameters := false

	for index, char := range path {
		if char == ':' {
			currentParameter = ""
			startGettingParameters = true
		} else if char == '/' && startGettingParameters {
			parameters = append(parameters, currentParameter)
			currentParameter = ""
			startGettingParameters = false
		} else if startGettingParameters {
			currentParameter += string(char)
		}

		if (index == len(path)-1) && startGettingParameters {
			parameters = append(parameters, currentParameter)
		}
	}

	return parameters
}
