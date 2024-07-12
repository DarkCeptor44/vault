package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	vaultErrors "github.com/DarkCeptor44/vault/internal/errors"
	"golang.org/x/crypto/argon2"
)

const (
	DocFolderName = "/documents"
	SaltLength    = 32
)

var (
	Docker    bool
	Debug     bool
	Host      string
	Port      int
	DocFolder string

	keyCache = NewCache()
)

func Init() {
	Docker = EnvBool("DOCKER")
	Host = Env("HOST", "0.0.0.0")
	Port = EnvInt("PORT", 8080)
	Debug = EnvBool("DEBUG")
	DocFolder = IfThenElse(Docker, DocFolderName, fmt.Sprint(".", DocFolderName))

	if Docker {
		log.Println("Running in Docker")
	}
	if Debug {
		log.Println("Debug mode")
	}

	_, err := os.Stat(DocFolder)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(DocFolder, os.ModePerm)
		HandleError(err, "Failed to create documents folder")
	}
}

func Env(name string, def string) string {
	aux := os.Getenv(name)
	if aux == "" {
		return def
	}
	return aux
}

func EnvBool(env string) bool {
	aux := os.Getenv(env)
	return strings.TrimSpace(aux) != ""
}

func EnvInt[V Integer](env string, def V) V {
	aux := os.Getenv(env)
	if aux == "" {
		return def
	}

	i, err := strconv.Atoi(aux)
	if err != nil {
		log.Println(env)
		return def
	}
	return V(i)
}

// If err is not nil and message is not empty, prints the message. If message is empty, prints the error
func HandleError(err error, message string) {
	if err != nil {
		if message != "" {
			log.Printf("%s: %s\n", message, err)
		} else {
			log.Println(err)
		}
		os.Exit(1)
	}
}

// Returns the first argument if the condition is true, the second if false.
//
// Original from https://github.com/shomali11/util/blob/f0771b70947f1d04b0e6826d333ab3ed295b05a0/xconditions/xconditions.go#L12. Modified for generics
func IfThenElse[V any](cond bool, a, b V) V {
	if cond {
		return a
	}
	return b
}

// Generates a random salt
func NewSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	return salt, err
}

// Derives a key from the password and salt using Argon2id, retrieves from the cache if already did before
func DeriveKey(pass, salt []byte) []byte {
	if key, ok := keyCache.Load(pass); ok {
		return key
	}

	key := argon2.IDKey(pass, salt, 3, 19*1024, 1, SaltLength)
	keyCache.Store(pass, key)
	return key
}

// Encrypts data with the given key using AES-256-GCM
func EncryptData(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypts data with the given key using AES-256-GCM
func DecryptData(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return nil, vaultErrors.ErrDataTooShort
	}

	nonce, data := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, data, nil)
}

// Removes spaces and converts filename to lowercase
func ClearFilename(filename string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(filename)), " ", "")
}
