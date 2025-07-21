package config

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

type Config struct {
	Port   string `env:"PORT"`
	Domain string `env:"DOMAIN"`
	Secret string `env:"SECRET"`
	DB     struct {
		DSN string `env:"DB_DSN"`
	}
	Dev bool `env:"DEV"`
}

func newConfig() (*Config, error) {
	resultCfg := &Config{}

	envCfg, err := envConfig()
	if err != nil {
		return nil, err
	}

	flagCfg, err := flagConfig()
	if err != nil {
		return nil, err
	}

	mergeConfig(resultCfg, flagCfg)
	mergeConfig(resultCfg, envCfg)

	err = checkParams(resultCfg)
	if err != nil {
		flag.Usage()
		return nil, err
	}

	return resultCfg, nil
}

func envConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func flagConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "p", "8080", "Service port")

	flag.StringVar(&cfg.Secret, "secrete", "", "Secret")

	flag.StringVar(&cfg.DB.DSN, "db-dsn", "", "PostgreSQL DSN")

	flag.BoolVar(&cfg.Dev, "dev", false, "Dev mode")

	flag.Parse()

	return cfg, nil
}

func mergeConfig(dst, src *Config) {
	if src == nil {
		return
	}

	if src.Port != "" {
		dst.Port = src.Port
	}
	if src.Secret != "" {
		dst.Secret = src.Secret
	}
	if src.DB.DSN != "" {
		dst.DB.DSN = src.DB.DSN
	}
	if src.Dev {
		dst.Dev = src.Dev
	}
}

func checkParams(cfg *Config) error {
	if cfg.Port != "" {
		port, err := strconv.Atoi(cfg.Port)
		if err != nil {
			return fmt.Errorf("port must be a number: %s", cfg.Port)
		}

		if port < 0 || port > 65535 {
			return fmt.Errorf("port must be between 0 and 65535 (inclusive): %d", port)
		}
	}

	if cfg.Secret == "" {
		return fmt.Errorf("secret must be specified")
	}

	switch cfg.Secret {
	case "":
		return fmt.Errorf("secret must be specified")
	default:
		if len(cfg.Secret) > 72 {
			return fmt.Errorf("secret length must be between 0 and 72, got %d", len(cfg.Secret))
		}
	}

	if cfg.DB.DSN == "" {
		return fmt.Errorf("PostgreSQL DSN must be specified")
	}

	return nil
}

func Provide() fx.Option {
	return fx.Provide(newConfig)
}
