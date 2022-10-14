package main

import (
	"archtecture/app"
	"archtecture/app/env"
	"archtecture/users"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	e, fileErr, osErr := env.NewEnv(".env")
	if fileErr != nil {
		fmt.Printf("%+v\n", fileErr)
	}
	if osErr != nil {
		fmt.Printf("%+v\n", osErr)
	}

	f := fiber.New()
	a := app.NewAppWithEvent(e, f)
	users.NewUserModule(a).Register()

	a.StartFiber()
}
