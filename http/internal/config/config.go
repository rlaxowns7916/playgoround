package config

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	DB struct {
		Database string `koanf:"database"`
		URL      string `koanf:"url"`
	} `koanf:"db"`
	HTTP struct {
		Port int `koanf:"port"`
	} `koanf:"http"`
}

func Load() *Config {
	config := &Config{}
	middlewares := []middleware{
		&EnvironmentValueReplaceMiddleware{},
	}
	phase := os.Getenv("APP_PHASE")
	k := koanf.New(".")
	err := k.Load(file.Provider("config/application.yaml"), yaml.Parser())

	if err := k.Load(file.Provider("config/application.yaml"), yaml.Parser()); err != nil {
		fmt.Printf("load base config error: %s", err)
		panic(err)
	}

	if phase != "" {
		if err := k.Load(file.Provider(fmt.Sprintf("config/application-%s.yaml", phase)), yaml.Parser()); err != nil {
			fmt.Printf("load phase config error: %s", err)
			panic(err)
		}
	}

	for _, key := range k.Keys() {
		value := k.Get(key)

		if strValue, ok := value.(string); ok {
			for _, mw := range middlewares {
				if err = mw.Handle(k, key, strValue); err != nil {
					fmt.Printf("middleware handle error: %s", err)
					panic(err)
				}
			}
		}
	}

	err = k.UnmarshalWithConf("", config, koanf.UnmarshalConf{Tag: "koanf"})
	if err != nil {
		fmt.Printf("unmarshal config error: %s\n", err)
		panic(err)
	}
	return config
}
