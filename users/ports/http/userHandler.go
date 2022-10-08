package http

import (
	"archtecture/app"
	"archtecture/app/validation"
	"archtecture/users/logic"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	logic *logic.User
}

func NewUserHandler(logic *logic.User) *UserHandler {
	return &UserHandler{logic: logic}
}

func (h *UserHandler) RegisterRoutes(a *app.App) {
	a.Fibre().Post("register", h.register())
}

func (h *UserHandler) register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := &logic.UserData{}
		err := ctx.BodyParser(input)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		user, err := h.logic.Register(input)
		if errorBag, ok := err.(*validation.ErrorBag); ok {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errorBag.All())
		}
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return ctx.Status(fiber.StatusCreated).JSON(user)
	}
}
