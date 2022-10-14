package http

import (
	"archtecture/app/http"
	"archtecture/users/logic"
	"github.com/gofiber/fiber/v2"
)

type loginInput struct {
	Identifier string `json:"identifier" form:"identifier"`
	Password   string `json:"password" form:"password"`
}

func (l *loginInput) to() {

}

type AuthHandler struct {
	authLogic *logic.Auth
	jwtAuth   *http.Auth
}

func NewAuthHandler(authLogic *logic.Auth, jwtAuth *http.Auth) *AuthHandler {
	return &AuthHandler{authLogic: authLogic, jwtAuth: jwtAuth}
}

func (h *AuthHandler) RegisterRoutes(f *fiber.App) {
	f.Post("/login", h.login())
	f.Post("/logout", h.logout(), http.MiddlewareAuthJwt())
}

func (h *AuthHandler) login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input loginInput

		if err := ctx.BodyParser(&input); err == fiber.ErrUnprocessableEntity {
			return ctx.SendStatus(fiber.StatusUnprocessableEntity)
		} else if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		inputMap := map[string]string{
			"identifier": input.Identifier,
			"password":   input.Password,
		}
		user, err := h.authLogic.AttemptLogin(inputMap)
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
