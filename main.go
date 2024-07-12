/**
 * Vault: Vault for encrypted text storage
 * Copyright (C) 2024 DarkCeptor44
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package main

import (
	"fmt"

	"github.com/DarkCeptor44/vault/internal/routes"
	"github.com/DarkCeptor44/vault/internal/util"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	util.Init()

	engine := html.New("./views", ".html")

	if util.Debug {
		engine.Reload(true)
	}

	app := fiber.New(fiber.Config{
		Views:       engine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Static("/static", "./public")

	app.Use(favicon.New())
	app.Use(idempotency.New())
	app.Use(recover.New())

	if util.Debug {
		app.Use(logger.New())
		app.Use(pprof.New())
	}

	// Routes
	routes.Route(app)

	err := app.Listen(fmt.Sprintf("%s:%d", util.Host, util.Port))
	util.HandleError(err, "")
}
