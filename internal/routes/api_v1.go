package routes

import (
	"errors"
	"os"
	"strings"

	vaultErrors "github.com/DarkCeptor44/vault/internal/errors"
	"github.com/DarkCeptor44/vault/internal/util"
	"github.com/gofiber/fiber/v2"
)

func routeApiV1(api fiber.Router) {
	v1 := api.Group("/v1")
	v1.Get("/check/:filename", func(c *fiber.Ctx) error {
		filename := util.ClearFilename(c.Params("filename"))

		dir, err := os.ReadDir(util.DocFolder)
		if err != nil {
			return fiberError(c, err, "")
		}

		if len(dir) == 0 {
			// no files exist, generate salt
			return notFound(c)
		}

		// theres files but we dont know if the right one exists
		for _, file := range dir {
			if file.IsDir() {
				continue
			}

			parts := strings.Split(file.Name(), "_")

			// shouldnt happen but just in case
			if len(parts) != 2 {
				continue
			}

			if strings.EqualFold(filename, parts[0]) {
				// file exists, retrieve salt
				return c.Status(fiber.StatusOK).JSON(&fiber.Map{
					"success": true,
					"message": "Salt is here",
					"salt":    parts[1],
				})
			}
		}

		return notFound(c)
	})

	v1.Post("/open", func(c *fiber.Ctx) error {
		var d struct {
			Filename string
			Hash     string
		}
		err := c.BodyParser(&d)
		if err != nil {
			return fiberError(c, err, "Could not parse body")
		}

		plainData, err := decrypt(util.ClearFilename(d.Filename), d.Hash)
		if errors.Is(err, vaultErrors.ErrInvalidKey) {
			return fiberError(c, err, "Invalid key provided", fiber.StatusUnauthorized)
		}
		if err != nil {
			return fiberError(c, err, "Could not decrypt data")
		}

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"message": "File decrypted",
			"data":    string(plainData),
		})
	})

	v1.Post("/save", func(c *fiber.Ctx) error {
		var data data
		err := c.BodyParser(&data)
		if err != nil {
			return fiberError(c, err, "Could not parse body")
		}

		data.Filename = util.ClearFilename(data.Filename)

		err = encrypt(data)
		if err != nil {
			return fiberError(c, err, "Could not encrypt data")
		}

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"message": "File saved",
		})
	})
}
