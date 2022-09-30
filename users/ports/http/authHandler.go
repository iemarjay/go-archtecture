package http

import (
	"archtecture/app"
	"archtecture/app/cache"
	"archtecture/app/http"
	"github.com/gofiber/fiber/v2"
)

type user interface {
	Key() string
}

type authLogic interface {
	AttemptLogin(interface{}) (user, error)
}

type loginRequest interface {
	toLogicInput() interface{}
}

type AuthHandler struct {
	authLogic authLogic
	cache     *cache.Cache
	jwtAuth   *http.Auth
}

func NewAuthHandler(authLogic authLogic, cache *cache.Cache) *AuthHandler {
	return &AuthHandler{authLogic: authLogic, cache: cache}
}

func (h *AuthHandler) RegisterRoutes(a *app.App) {
	a.Fibre().Post("login", h.login())
}

func (h *AuthHandler) login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input loginRequest
		user, err := h.authLogic.AttemptLogin(input.toLogicInput())
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		var token string
		token, err = h.jwtAuth.CreateToken(user.Key())
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.JSON(fiber.Map{"token": token, "user": user})
	}
}

func (h *AuthHandler) logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := h.jwtAuth.DestroyToken(ctx)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
