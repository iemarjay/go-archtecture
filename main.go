package main

import (
	"archtecture/app"
	"archtecture/app/env"
	"fmt"
)

func main() {
	e, err := env.NewEnv()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	a := app.NewApp(e)

	a.StartFiber()
}
