package main

import (
	"archtecture/app"
	"archtecture/app/env"
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
	a := app.NewApp(e, f)

	a.StartFiber()
}
