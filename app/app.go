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
	fiber *fiber.App
}

func NewApp(env *env.Env, fibre *fiber.App) *App {
	return &App{env: env, fiber: fibre}
}

func (a *App) Env() *env.Env {
	return a.env
}

func (a *App) Database() *database.MongoDatabase {
	config := database.NewMongoDatabaseConfigFromEnv(a.env)
	return database.NewMongoDatabaseFromConfig(config)
}

func (a *App) Cache() *cache.Redis {
	return cache.NewCacheWithRedisFromConfig(cache.NewConfig(a.env))
}

func (a *App) StartFiber() {
	a.fiber.Static(a.env.PublicPathPrefix, a.env.PublicRootDir)
	a.fiber.Use(cors.New())

	log.Fatal(a.fiber.Listen(":" + a.env.Port))
}

func (a *App) Fiber() *fiber.App {
	return a.fiber
}
