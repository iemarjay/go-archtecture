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

func NewEnv(path ...string) (*Env, error, error) {
	var fileEnvMap map[string]string
	var err error
	if len(path) > 0 {
		fileEnvMap, err = godotenv.Read(path[0])
	}

	env, err2 := NewEnvFromMap(fileEnvMap)

	return env, err, err2
}

func NewEnvFromMap(defaults map[string]string) (*Env, error) {
	e := &Env{}
	options := eEnv.Options{}

	if len(defaults) > 0 {
		options = eEnv.Options{
			Environment: defaults,
		}
	}

	if err := eEnv.Parse(e, options); err != nil {
		return nil, err
	}

	return e, nil
}
