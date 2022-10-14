package env

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_EnvCanBeLoadedFromMap(t *testing.T) {
	defaults := map[string]string{
		"SMS_FROM": "Arch",
		"MailFrom": "Today",
	}
	env, _ := NewEnvFromMap(defaults)

	assert.Equal(t, env.SmsFrom, defaults["SMS_FROM"])
	assert.Equal(t, env.MailFrom, defaults["MailFrom"])
}

func Test_EnvCanBeLoadedFromFromEnvFile(t *testing.T) {
	env, _, _ := NewEnv("test.env")

	assert.Equal(t, env.SmsFrom, "Arch")
	assert.Equal(t, env.MailFrom, "Today")
}
