package config

import (
	"time"

	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
	oauth2google "golang.org/x/oauth2/google"
)

type (
	Config struct {
		App         `yaml:"app"`
		HTTP        `yaml:"http"`
		Log         `yaml:"logger"`
		PG          `yaml:"postgres"`
		MongoDB     `yaml:"mongodb"`
		Cache       `yaml:"cache"`
		Redis       `yaml:"redis"`
		SocialAuth  `yaml:"social_auth"`
		Session     `yaml:"session"`
		AccessToken `yaml:"access_token"`
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

	SocialAuth struct {
		GitHubClientID     string `yaml:"github_client_id" env-required:"true" env:"GH_CLIENT_ID"`
		GitHubClientSecret string `env-required:"true" env:"GH_CLIENT_SECRET"`
		GitHubScope        string `yaml:"github_scope" env-required:"true" env:"GH_SCOPE"`

		GoogleClientID     string `yaml:"google_client_id" env-required:"true" env:"GOOGLE_CLIENT_ID"`
		GoogleClientSecret string `env-required:"true" env:"GOOGLE_CLIENT_SECRET"`
		GoogleScope        string `yaml:"google_scope" env-required:"true" env:"GOOGLE_SCOPE"`
	}

	Session struct {
		TTL            time.Duration `env-required:"true" yaml:"ttl" env:"SESSION_TTL"`
		CookieKey      string        `env-required:"true" yaml:"cookie_key" env:"SESSION_COOKIE_KEY"`
		CookieDomain   string        `yaml:"cookie_domain" env:"SESSION_COOKIE_DOMAIN"`
		CookieSecure   bool          `yaml:"cookie_secure" env:"SESSION_COOKIE_SECURE"`
		CookieHTTPOnly bool          `yaml:"cookie_httponly" env:"SESSION_COOKIE_HTTPONLY"`
	}

	AccessToken struct {
		TTL        time.Duration `env-required:"true" yaml:"ttl" env:"ACCESS_TOKEN_TTL"`
		SigningKey string        `env-required:"true" yaml:"signing_key" env:"ACCESS_TOKEN_SIGNING_KEY"`
	}
)

func (sa *SocialAuth) Endpoints() map[string]oauth2.Endpoint {
	return map[string]oauth2.Endpoint{
		"github": oauth2github.Endpoint,
		"google": oauth2google.Endpoint,
	}
}

func (sa *SocialAuth) Scopes() map[string]string {
	return map[string]string{
		"github": sa.GitHubScope,
		"google": sa.GoogleScope,
	}
}

func (sa *SocialAuth) ClientIDs() map[string]string {
	return map[string]string{
		"github": sa.GitHubClientID,
		"google": sa.GoogleClientID,
	}

}

func (sa *SocialAuth) ClientSecrets() map[string]string {
	return map[string]string{
		"github": sa.GitHubClientSecret,
		"google": sa.GoogleClientSecret,
	}

}
