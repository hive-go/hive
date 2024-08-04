package hive

import (
	"github.com/gofiber/fiber/v2"
)

type Map = fiber.Map

type Ctx = fiber.Ctx

type Response interface {
	Map | string
}
