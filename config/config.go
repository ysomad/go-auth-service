package config

import "time"

type (
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		PG    `yaml:"postgres"`
		Redis `yaml:"redis"`
		JWT   `yaml:"jwt"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" env:"PG_URL"`
	}

	Redis struct {
		Addr     string `env-required:"true" env:"REDIS_ADDR"`
		Password string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	JWT struct {
		AccessTokenTTL  time.Duration `env-required:"true" yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL"`
		RefreshTokenTTL time.Duration `env-required:"true" yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL"`
		SigningKey      string        `env-required:"true" yaml:"signing_key" env:"SIGNING_KEY"`
	}
)
