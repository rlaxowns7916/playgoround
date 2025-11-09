package config

import (
	"os"
	"regexp"

	"github.com/knadh/koanf/v2"
)

type middleware interface {
	Handle(k *koanf.Koanf, key, value string) error
}

type EnvironmentValueReplaceMiddleware struct{}

func (e *EnvironmentValueReplaceMiddleware) Handle(k *koanf.Koanf, key, value string) error {
	replaced := e.replace(value)
	return k.Set(key, replaced)
}

func (e *EnvironmentValueReplaceMiddleware) replace(value string) string {
	re := regexp.MustCompile(`\$\{([^}:]+)(?::([^}]*))?\}`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		submatch := re.FindStringSubmatch(match)
		varName := submatch[1]

		if val, exists := os.LookupEnv(varName); exists {
			return val
		}

		return match
	})
}
