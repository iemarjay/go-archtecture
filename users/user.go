package users

import (
	"archtecture/app"
	appHttp "archtecture/app/http"
	"archtecture/app/notification"
	"archtecture/app/notification/channels"
	"archtecture/app/validation"
	"archtecture/users/logic"
	"archtecture/users/logic/messages"
	"archtecture/users/ports/http"
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

func (u *UserModule) Register() {
	u.app.Fibre().Use(appHttp.MiddlewareAuthUser(u.makeMongoRepository(), u.app.Cache()))

	http.NewAuthHandler(u.makeAuthLogic(), u.makeJwtAuth()).RegisterRoutes(u.app)
	http.NewUserHandler(u.makeUserLogic()).RegisterRoutes(u.app)
}

func (u *UserModule) makeAuthLogic() *logic.Auth {
	return logic.NewAuth(u.makeMongoRepository())
}

func (u *UserModule) makeUserLogic() *logic.User {
	env := u.app.Env()
	repository := u.makeMongoRepository()
	validator := validation.NewValidator()
	notifier := notification.NewDefaultNotifier()
	sms := channels.NewSmsFromEnv(env)
	email := channels.NewMailgunFromEnv(env)
	message := messages.NewWelcome(sms, email)

	return logic.NewUser(repository, validator, notifier, message)
}

func (u *UserModule) makeMongoRepository() *repositories.Mongo {
	database := u.app.Database()
	database.Table(repositories.TableName)

	return repositories.NewMongo(database)
}

func (u *UserModule) makeJwtAuth() *appHttp.Auth {
	return appHttp.NewAuth(u.makeMongoRepository(), u.app.Cache())
}
