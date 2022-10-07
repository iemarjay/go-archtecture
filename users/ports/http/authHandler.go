package http

import (
	"archtecture/app"
	"archtecture/app/cache"
	"archtecture/app/http"
	"archtecture/users/logic"
	"github.com/gofiber/fiber/v2"
)

type authLogic interface {
	AttemptLogin(map[string]string) (*logic.UserData, error)
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
	a.Fibre().Get("", func(ctx *fiber.Ctx) error {
		return ctx.JSON("message")
	})
	a.Fibre().Post("/login", h.login())
	a.Fibre().Post("/logout", h.logout(), http.MiddlewareAuthJwt())
}

func (h *AuthHandler) login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := map[string]string{
			"identifier": ctx.FormValue("identifier"),
			"password":   ctx.FormValue("password"),
		}

		user, err := h.authLogic.AttemptLogin(input)
		if err == logic.UserNotFound {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		} else if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
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
