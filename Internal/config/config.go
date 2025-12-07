package config

import (
	"errors"
	"log"
	"net/url"

	"github.com/Kenasvarghese/Caching-Proxy/Internal/proxy"
	"github.com/Kenasvarghese/Caching-Proxy/Internal/rate_limiter"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Origin    string `default:"http://localhost" split_words:"true"`
	Port      int    `default:"8080" split_words:"true"`
	originUrl *url.URL

	rate_limiter.RateLimiterConfig

	proxy.TransportConfig
}

// LoadConfig loads configuration from environment variables into a Config struct.
func LoadConfig() *Config {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("error loading config:", err)
	}

	// validate only critical values
	if err := cfg.Validate(); err != nil {
		log.Fatal("invalid config:", err)
	}

	return &cfg
}

// Validate validates the config struct
func (c *Config) Validate() error {
	originURL, err := url.Parse(c.Origin)
	if err != nil {
		return errors.New("invalid origin url")
	}
	c.originUrl = originURL
	return nil
}
func (c *Config) GetOriginURL() *url.URL {
	return c.originUrl
}
