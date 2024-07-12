package routes

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	vaultErrors "github.com/DarkCeptor44/vault/internal/errors"
	"github.com/DarkCeptor44/vault/internal/util"
	"github.com/gofiber/fiber/v2"
)

func encrypt(data data) error {
	newSalt, err := hex.DecodeString(data.Salt)
	if err != nil {
		return err
	}

	key := util.DeriveKey([]byte(data.Hash), []byte(data.Salt))
	cipherText, err := util.EncryptData(key, []byte(data.Text))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(util.DocFolder, fmt.Sprintf("%s_%s", data.Filename, hex.EncodeToString(newSalt))), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()

	_, err = bw.Write(cipherText)
	return err
}

func decrypt(filename, hash string) (plainData []byte, err error) {
	salt, fullname, err := saltFromFile(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(util.DocFolder, fullname))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, errors.New("file is empty")
	}

	key := util.DeriveKey([]byte(hash), []byte(salt))
	plainData, err = util.DecryptData(key, bytes)
	if err != nil {
		return nil, vaultErrors.ErrInvalidKey
	}

	return plainData, nil
}

func fiberError(c *fiber.Ctx, err error, message string, code ...int) error {
	log.Println(err)
	return c.Status(util.IfThenElse(len(code) > 0, code[0], fiber.StatusInternalServerError)).JSON(&fiber.Map{
		"success": false,
		"message": util.IfThenElse(message != "", message, err.Error()),
	})
}

func notFound(c *fiber.Ctx) error {
	salt, err := util.NewSalt()
	if err != nil {
		return fiberError(c, err, "Could not generate salt")
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": false,
		"message": "File not found",
		"salt":    hex.EncodeToString(salt),
	})
}

func saltFromFile(filename string) (salt, fullname string, err error) {
	dir, err := os.ReadDir(util.DocFolder)
	if err != nil {
		return "", "", err
	}

	if len(dir) == 0 {
		return "", "", errors.New("file not found")
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		parts := strings.Split(file.Name(), "_")
		if len(parts) != 2 {
			continue
		}

		if strings.EqualFold(filename, parts[0]) {
			return parts[1], file.Name(), nil
		}
	}

	return "", "", errors.New("file not found")
}
