package users

import (
	"archtecture/app"
	appHttp "archtecture/app/http"
	"archtecture/users/ports/http"
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
	http.NewAuthHandler(authLogic, cache).RegisterRoutes(u.app)

	u.app.Fibre().Use(appHttp.MiddlewareAuthUser(db, cache))
}
