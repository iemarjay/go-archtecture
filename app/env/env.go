package env

import (
	eEnv "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Env struct {
	Port string `env:"PORT" envDefault:"3000"`

	DatabaseUrl  string `env:"DATABASE_URL"`
	DatabaseName string `env:"DATABASE_NAME"`

	PublicPathPrefix string `env:"PUBLIC_PATH_PREFIX" envDefault:"public"`
	PublicRootDir    string `env:"PUBLIC_ROOT_DIR" envDefault:"resources/public"`
	RedisUrl         string `env:"REDIS_URL" envDefault:"resources/public"`
	RedisPassword    string `env:"REDIS_PASSWORD"`

	MailgunDomain     string `env:"MAILGUN_DOMAIN"`
	MailgunPrivateKey string `env:"MAILGUN_PRIVATE_KEY"`
	MailFrom          string `env:"MailFrom" envDefault:"mail@example.com"`

	SmsFrom         string `env:"SMS_FROM" envDefault:"architecture"`
	TermiiUri       string `env:"TERMII_URI" envDefault:"https://api.ng.termii.com/api/sms/send"`
	TermiiApiKey    string `env:"TERMII_API_KEY"`
	InfobipUri      string `env:"INFOBIP_URI"`
	InfobipUsername string `env:"INFOBIP_USERNAME"`
	InfobipPassword string `env:"INFOBIP_PASSWORD"`
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
