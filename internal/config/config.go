package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"time"
)

const (
	CONFIG_DIR       = "configs"
	CONFIG_PROD_FILE = "prod"
	CONFIG_DEV_FILE  = "dev"
)

type Config struct {
	Postgres Postgres
	Url      Url
	Limiter  Limiter `mapstructure:"limiter"`
	Auth     Auth    `mapstructure:"auth"`
	IsProd   bool

	Http struct {
		Port int64 `mapstructure:"port"`
	} `mapstructure:"http"`

	Ctx struct {
		Ttl time.Duration `mapstructure:"ttl"`
	} `mapstructure:"ctx"`
}

type Postgres struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

type Url struct {
	Warehouse string
	Accounts  string
}

type Limiter struct {
	RPS   int           `mapstructure:"rps"`
	Burst int           `mapstructure:"burst"`
	TTL   time.Duration `mapstructure:"ttl"`
}

type Auth struct {
	AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	SigningKey      string        `vault:"auth_signing_key"`
}

func New(isProd bool) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(CONFIG_DIR)
	viper.SetConfigName(CONFIG_DEV_FILE)

	if isProd {
		viper.SetConfigName(CONFIG_PROD_FILE)
		cfg.IsProd = true
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("postgres", &cfg.Postgres); err != nil {
		return nil, err
	}

	if err := envconfig.Process("url", &cfg.Url); err != nil {
		return nil, err
	}

	if err := envconfig.Process("auth", &cfg.Auth); err != nil {
		return nil, err
	}

	return cfg, nil
}
