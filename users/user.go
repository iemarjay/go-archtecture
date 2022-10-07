package users

import (
	"archtecture/app"
	appHttp "archtecture/app/http"
	"archtecture/users/logic"
	"archtecture/users/ports/http"
	"archtecture/users/repository"
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
	cache := u.app.Cache()
	u.app.Fibre().Use(appHttp.MiddlewareAuthUser(u.makeMongoRepository(), cache))

	http.NewAuthHandler(u.makeAuthLogic(), cache).RegisterRoutes(u.app)
}

func (u *UserModule) makeAuthLogic() *logic.Auth {
	return logic.NewAuth(u.makeMongoRepository())
}

func (u *UserModule) makeMongoRepository() *repository.Mongo {
	database := u.app.Database()
	database.Table(repository.TableName)

	return repository.NewMongo(database)
}
