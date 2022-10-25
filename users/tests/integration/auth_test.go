package integration

import (
	"archtecture/app/cache"
	appHttp "archtecture/app/http"
	"archtecture/users/logic"
	"archtecture/users/ports/rest"
	"archtecture/users/repositories"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin_Successful(t *testing.T) {
	f := fiber.New()
	user := &logic.UserData{
		ID:        "",
		Lastname:  "",
		Firstname: "",
		Phone:     "+2349083874378",
		Email:     "auth@example.test",
		Password:  "secret",
	}
	_ = user.EncryptPassword()
	repository := repositories.NewMap(user)
	rest.NewAuthHandler(logic.NewAuth(repository), makeJwtAuth(repository)).
		RegisterRoutes(f)

	input := `{"identifier": "auth@example.test", "password": "secret"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := f.Test(req, -1)
	body, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	u := data["user"].(map[string]interface{})
	token, exists := data["token"]

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)
	assert.Equal(t, true, exists)
	assert.NotEqual(t, "", token)
	assert.Equal(t, "auth@example.test", u["email"])
	assert.Equal(t, "+2349083874378", u["phone"])
}

func TestLogin_ValidationFails(t *testing.T) {
	f := fiber.New()
	user := &logic.UserData{
		ID:        "",
		Lastname:  "",
		Firstname: "",
		Phone:     "+2349083874378",
		Email:     "auth@example.test",
		Password:  "secret",
	}
	_ = user.EncryptPassword()

	repository := repositories.NewMap(user)
	rest.NewAuthHandler(logic.NewAuth(repository), makeJwtAuth(repository)).
		RegisterRoutes(f)

	input := `{"identifier": "auth@example.test", "password": "secret3234"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := f.Test(req, -1)

	assert.Equal(t, resp.StatusCode, fiber.StatusUnauthorized)
}

func makeJwtAuth(repo *repositories.Map) *appHttp.Auth {
	return appHttp.NewAuth(repo, mapCache())
}

func mapCache() *cache.Map {
	return cache.NewMap(map[string]*cache.Data{})
}

type httpTest struct {
}
