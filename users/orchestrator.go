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

type UserModule struct {
	app *app.App
}

func NewUserModule(a *app.App) *UserModule {
	return &UserModule{
		app: a,
	}
}

func (u *UserModule) BootWithMongoAndFiber() {
	u.app.Fiber().Use(appHttp.MiddlewareAuthUser(u.makeMongoRepository(), u.app.Cache()))

	rest.NewAuthHandler(u.makeAuthLogic(), u.makeJwtAuth()).RegisterRoutes(u.app.Fiber())
	rest.NewUserHandler(u.makeUserLogic()).RegisterRoutes(u.app.Fiber())

	u.app.Event().Listen(logic.UserRegisteredEvent, u.makeSendWelcomeMail())
}

func (u *UserModule) makeAuthLogic() *logic.Auth {
	return logic.NewAuth(u.makeMongoRepository())
}

func (u *UserModule) makeUserLogic() *logic.User {
	repository := u.makeMongoRepository()
	validator := validation.NewValidator()

	return logic.NewUser(repository, validator, u.app.Event())
}

func (u *UserModule) makeMongoRepository() *repositories.Mongo {
	database := u.app.Database()
	database.Table(repositories.TableName)

	return repositories.NewMongo(database)
}

func (u *UserModule) makeJwtAuth() *appHttp.Auth {
	return appHttp.NewAuth(u.makeMongoRepository(), u.app.Cache())
}

func (u *UserModule) makeSendWelcomeMail() *listeners.SendWelcomeMail {
	mg := channels.NewMailgunFromEnv(u.app.Env())
	return listeners.NewSendWelcomeMail(mg)
}
