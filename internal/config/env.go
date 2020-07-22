package config

import (
	"os"
)

func GetEnv(name, def string) string {
	variable := os.Getenv(name)
	if variable != "" {
		return variable
	}
	return def
}
