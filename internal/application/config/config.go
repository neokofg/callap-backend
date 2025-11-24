package config

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env            Env    `env:"ENV"             env-required:"true"`
	ContextTimeout int    `env:"CONTEXT_TIMEOUT" env-required:"true"`
	Host           string `env:"HOST"            env-required:"true"`
	Port           string `env:"PORT"            env-required:"true"`
	JWT            JWT    `                      env-required:"true" env-prefix:"JWT_"`
}

type JWT struct {
	Secret          string `env:"SECRET"            env-required:"true"`
	AccessTokenTTL  int    `env:"ACCESS_TOKEN_TTL"  env-required:"true"`
	RefreshTokenTTL int    `env:"REFRESH_TOKEN_TTL" env-required:"true"`
	Issuer          string `env:"ISSUER"            env-default:"backend"`
}
