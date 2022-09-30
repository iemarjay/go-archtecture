package env

import (
	eEnv "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Env struct {
	Port string `env:"PORT" envDefault:"3000"`

	DatabaseUrl  string `env:"DATABASE_URL"`
	DatabaseName string `env:"DATABASE_NAME"`

	TermiiKey string `env:"TERMII_KEY"`

	PublicPathPrefix string `env:"PUBLIC_PATH_PREFIX" envDefault:"public"`
	PublicRootDir    string `env:"PUBLIC_ROOT_DIR" envDefault:"resources/public"`
	RedisUrl         string `env:"REDIS_URL" envDefault:"resources/public"`
	RedisPassword    string `env:"REDIS_PASSWORD"`
}

func NewEnv() (*Env, error) {
	e := &Env{}
	if err := eEnv.Parse(e); err != nil {
		return nil, err
	}

	return e, nil
}

func NewEnvFromMap(defaults map[string]string) (*Env, error) {
	e := &Env{}
	options := eEnv.Options{
		Environment: defaults,
	}
	if err := eEnv.Parse(e, options); err != nil {
		return nil, err
	}

	return e, nil
}

func NewEnvFromFile(filePath string) (*Env, error) {
	value, err := godotenv.Read(filePath)
	if err != nil {
		return nil, err
	}
	env, err := NewEnvFromMap(value)
	if err != nil {
		return nil, err
	}

	return env, err
}
