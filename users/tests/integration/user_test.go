package integration

import (
	"archtecture/app/events"
	"archtecture/app/validation"
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

func TestRegisterEndpoint_IsOk(t *testing.T) {
	f := fiber.New()
	user := &logic.UserData{
		Lastname:  "Paul",
		Firstname: "Dream",
		Phone:     "+2349083874378",
		Email:     "auth@example.test",
		Password:  "secret",
	}
	_ = user.EncryptPassword()

	repository := repositories.NewMap()
	rest.NewUserHandler(logic.NewUser(repository, validation.NewValidator(), events.NewEvent())).
		RegisterRoutes(f)

	input := `{"email": "auth@example.test", "phone": "+2349083874378", "password": "secret", "firstname": "Dream", "lastname": "Paul"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := f.Test(req, -1)
	body, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	assert.Equal(t, "auth@example.test", data["email"])
	assert.Equal(t, "+2349083874378", data["phone"])

	u, _ := repository.FindByUniqueFields("auth@example.test")
	assert.NotEqual(t, u, nil)
}

func TestRegisterEndpoint_ShouldRespondWithUnprocessableEntity(t *testing.T) {
	f := fiber.New()

	repository := repositories.NewMap()
	rest.NewUserHandler(logic.NewUser(repository, validation.NewValidator(), events.NewEvent())).
		RegisterRoutes(f)

	input := `{"password": "secret", "firstname": "Dream", "lastname": "Paul"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := f.Test(req, -1)
	body, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	assert.Equal(t, resp.StatusCode, fiber.StatusUnprocessableEntity)
	_, exists := data["UserData.Email"]
	assert.Equal(t, exists, true)
	u, _ := repository.FindByUniqueFields("auth@example.test")
	assert.Equal(t, u, nil)
}
