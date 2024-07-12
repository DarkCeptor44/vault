package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	routeRoot(app)

	api := app.Group("/api")
	routeApiV1(api)

	app.Use(func(c *fiber.Ctx) error {
		return c.Redirect("/404")
	})
}

func routeRoot(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil, "layouts/main")
	})

	app.All("/404", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", &fiber.Map{
			"Title": "Page Not Found",
		}, "layouts/main")
	})
}
