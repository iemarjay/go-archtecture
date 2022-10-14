package app

import (
	"archtecture/app/cache"
	"archtecture/app/database"
	"archtecture/app/env"
	"archtecture/app/events"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

type App struct {
	env   *env.Env
	event *events.Event
	fiber *fiber.App
}

func NewAppWithEvent(env *env.Env, fiber *fiber.App) *App {
	return NewApp(env, fiber, events.NewEvent())
}

func NewApp(env *env.Env, fibre *fiber.App, event *events.Event) *App {
	return &App{env: env, fiber: fibre, event: event}
}

func (a *App) Env() *env.Env {
	return a.env
}

func (a *App) Event() *events.Event {
	return a.event
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
