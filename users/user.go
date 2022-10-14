package users

import (
	"archtecture/app"
	appHttp "archtecture/app/http"
	"archtecture/app/notification"
	"archtecture/app/notification/channels"
	"archtecture/app/validation"
	"archtecture/users/listeners"
	"archtecture/users/logic"
	"archtecture/users/messages"
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
	u.app.Fiber().Use(appHttp.MiddlewareAuthUser(u.makeMongoRepository(), u.app.Cache()))

	http.NewAuthHandler(u.makeAuthLogic(), u.makeJwtAuth()).RegisterRoutes(u.app.Fiber())
	http.NewUserHandler(u.makeUserLogic()).RegisterRoutes(u.app)

	u.app.Event().Listen(logic.UserRegisteredName, u.makeSendWelcomeMessageListener())
}

func (u *UserModule) makeAuthLogic() *logic.Auth {
	return logic.NewAuth(u.makeMongoRepository())
}

func (u *UserModule) makeUserLogic() *logic.User {
	repository := u.makeMongoRepository()
	validator := validation.NewValidator()

	return logic.NewUser(repository, validator, u.app.Event())
}

func (u *UserModule) makeSendWelcomeMessageListener() *listeners.SendWelcomeNotification {
	notifier := notification.NewDefaultNotifier()
	message := u.makeWelcomeMessage()

	return listeners.NewSendWelcomeNotification(notifier, message)
}

func (u *UserModule) makeWelcomeMessage() *messages.Welcome {
	env := u.app.Env()
	sms := channels.NewSmsFromEnv(env)
	email := channels.NewMailgunFromEnv(env)
	return messages.NewWelcome(sms, email)
}

func (u *UserModule) makeMongoRepository() *repositories.Mongo {
	database := u.app.Database()
	database.Table(repositories.TableName)

	return repositories.NewMongo(database)
}

func (u *UserModule) makeJwtAuth() *appHttp.Auth {
	return appHttp.NewAuth(u.makeMongoRepository(), u.app.Cache())
}
