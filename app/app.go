package app

import (
	"archtecture/app/cache"
	"archtecture/app/database"
	"archtecture/app/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

type App struct {
	env   *env.Env
	fibre *fiber.App
}

func NewApp(env *env.Env, fibre *fiber.App) *App {
	return &App{env: env, fibre: fibre}
}

func (a *App) Env() *env.Env {
	return a.env
}

func (a *App) Database() *database.MongoDatabase {
	config := database.NewMongoDatabaseConfigFromEnv(a.env)
	return database.NewMongoDatabaseFromConfig(config)
}

func (a *App) Cache() *cache.Cache {
	return cache.NewCacheWithRedisFromConfig(cache.NewConfig(a.env))
}

func (a *App) StartFiber() {
	a.fibre.Static(a.env.PublicPathPrefix, a.env.PublicRootDir)
	a.fibre.Use(cors.New())

	log.Fatal(a.fibre.Listen(":" + a.env.Port))
}

func (a *App) Fibre() *fiber.App {
	return a.fibre
}
