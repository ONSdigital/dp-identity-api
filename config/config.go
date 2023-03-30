package config

import (
	"encoding/json"
	"time"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-identity-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	AWSRegion                  string        `envconfig:"AWS_REGION"`
	AWSCognitoUserPoolID       string        `envconfig:"AWS_COGNITO_USER_POOL_ID" json:"-"`
	AWSCognitoClientId         string        `envconfig:"AWS_COGNITO_CLIENT_ID" json:"-"`
	AWSCognitoClientSecret     string        `envconfig:"AWS_COGNITO_CLIENT_SECRET" json:"-"`
	AWSAuthFlow                string        `envconfig:"AWS_AUTH_FLOW" json:"-"`
	AllowedEmailDomains        []string      `envconfig:"ALLOWED_EMAIL_DOMAINS" json:"-"`
	AuthorisationConfig        *authorisation.Config
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25600",
		GracefulShutdownTimeout:    20 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		AWSRegion:                  "eu-west-2",
		AWSAuthFlow:                "USER_PASSWORD_AUTH",
		AllowedEmailDomains:        []string{"@ons.gov.uk", "@ext.ons.gov.uk"},
		AuthorisationConfig:        authorisation.NewDefaultConfig(),
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Config) String() string {
	configJson, _ := json.Marshal(config)
	return string(configJson)
}
