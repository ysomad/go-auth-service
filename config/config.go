package config

import "time"

type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"logger"`
		PG      `yaml:"postgres"`
		MongoDB `yaml:"mongodb"`
		Cache   `yaml:"cache"`
		Redis   `yaml:"redis"`
		Session `yaml:"session"`
		JWT     `yaml:"jwt"`
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

	MongoDB struct {
		URI      string `env-required:"true" env:"MONGO_URI"`
		Username string `env-required:"true" env:"MONGO_USER"`
		Password string `env-required:"true" env:"MONGO_PASS"`
		Database string `env-required:"true" yaml:"database" env:"MONGO_DATABASE"`
	}

	Cache struct {
		TTL time.Duration `env-required:"true" yaml:"ttl" env:"CACHE_TTL"`
	}

	Redis struct {
		Addr     string `env-required:"true" env:"REDIS_ADDR"`
		Password string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	Session struct {
		TTL            time.Duration `env-required:"true" yaml:"ttl" env:"SESSION_TTL"`
		CookieKey      string        `env-required:"true" yaml:"cookie_key" env:"SESSION_COOKIE_KEY"`
		CookieDomain   string        `yaml:"cookie_domain" env:"SESSION_COOKIE_DOMAIN"`
		CookieSecure   bool          `yaml:"cookie_secure" env:"SESSION_COOKIE_SECURE"`
		CookieHttpOnly bool          `yaml:"cookie_httponly" env:"SESSION_COOKIE_HTTPONLY"`
	}

	JWT struct {
		TTL        time.Duration `env-required:"true" yaml:"ttl" env:"JWT_TTL"`
		SigningKey string        `env-required:"true" yaml:"signing_key" env:"SIGNING_KEY"`
	}
)
