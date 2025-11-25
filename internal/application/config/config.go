package config

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env            Env        `env:"ENV"             env-required:"true"`
	ContextTimeout int        `env:"CONTEXT_TIMEOUT" env-required:"true"`
	Host           string     `env:"HOST"            env-required:"true"`
	Port           string     `env:"PORT"            env-required:"true"`
	JWT            JWT        `                      env-required:"true" env-prefix:"JWT_"`
	PostgreSQL     PostgreSQL `                      env-required:"true" env-prefix:"POSTGRES_"`
}

type JWT struct {
	Secret          string `env:"SECRET"            env-required:"true"`
	AccessTokenTTL  int    `env:"ACCESS_TOKEN_TTL"  env-required:"true"`
	RefreshTokenTTL int    `env:"REFRESH_TOKEN_TTL" env-required:"true"`
	Issuer          string `env:"ISSUER"            env-default:"backend"`
}

type PostgreSQL struct {
	Username string         `env:"USERNAME" env-required:"true"`
	Password string         `env:"PASSWORD" env-required:"true"`
	Host     string         `env:"HOST"     env-required:"true"`
	Port     string         `env:"PORT"     env-required:"true"`
	Database string         `env:"DATABASE" env-required:"true"`
	Pool     PostgreSQLPool `               env-required:"true" env-prefix:"POOL_"`
}

type PostgreSQLPool struct {
	MaxConns          int `env:"MAX_CONNS"           env-required:"true"`
	MinConns          int `env:"MIN_CONNS"           env-required:"true"`
	MaxConnLifeTime   int `env:"MAX_CONN_LIFE_TIME"  env-required:"true"`
	MaxConnIdleTime   int `env:"MAX_CONN_IDLE_TIME"  env-required:"true"`
	HealthCheckPeriod int `env:"HEALTH_CHECK_PERIOD" env-required:"true"`
}
