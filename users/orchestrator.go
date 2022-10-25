package users

import (
	"archtecture/app"
	appHttp "archtecture/app/http"
	"archtecture/app/notification/channels"
	"archtecture/app/validation"
	"archtecture/users/listeners"
	"archtecture/users/logic"
	"archtecture/users/ports/rest"
	"archtecture/users/repositories"
)

type Orchestrator struct {
	app *app.App
}

func NewOrchestrator(a *app.App) *Orchestrator {
	return &Orchestrator{
		app: a,
	}
}

func (u *Orchestrator) BootWithMongoAndFiber() {
	u.app.Fiber().Use(appHttp.MiddlewareAuthUser(u.makeMongoRepository(), u.app.Cache()))

	rest.NewAuthHandler(u.makeAuthLogic(), u.makeJwtAuth()).RegisterRoutes(u.app.Fiber())
	rest.NewUserHandler(u.makeUserLogic()).RegisterRoutes(u.app.Fiber())

	u.app.Event().Listen(logic.UserRegisteredEvent, u.makeSendWelcomeMail())
}

func (u *Orchestrator) makeAuthLogic() *logic.Auth {
	return logic.NewAuth(u.makeMongoRepository())
}

func (u *Orchestrator) makeUserLogic() *logic.User {
	repository := u.makeMongoRepository()
	validator := validation.NewValidator()

	return logic.NewUser(repository, validator, u.app.Event())
}

func (u *Orchestrator) makeMongoRepository() *repositories.Mongo {
	database := u.app.Database()
	database.Table(repositories.TableName)

	return repositories.NewMongo(database)
}

func (u *Orchestrator) makeJwtAuth() *appHttp.Auth {
	return appHttp.NewAuth(u.makeMongoRepository(), u.app.Cache())
}

func (u *Orchestrator) makeSendWelcomeMail() *listeners.SendWelcomeMail {
	mg := channels.NewMailgunFromEnv(u.app.Env())
	return listeners.NewSendWelcomeMail(mg)
}
