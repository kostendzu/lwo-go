package config

import (
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
)

func Searchup(dir string, filename string) string {
	if dir == "/" || dir == "" {
		return ""
	}

	if _, err := os.Stat(path.Join(dir, filename)); err == nil {
		return path.Join(dir, filename)
	}

	return Searchup(path.Dir(dir), filename)
}

func Load(path string) error {
	directory, err := os.Getwd()
	if err != nil {
		return err
	}

	filepath := Searchup(directory, ".env")

	if filepath == "" {
		return fmt.Errorf("could not find dot env file")
	}

	return godotenv.Load(filepath)
}
