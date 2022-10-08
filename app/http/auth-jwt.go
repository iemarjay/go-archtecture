package http

import (
	"archtecture/app/cache"
	"archtecture/app/utils"
	"archtecture/users/logic"
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type db interface {
	Find(id string) (*logic.UserData, error)
}

func MiddlewareAuthJwt() fiber.Handler {
	return jwtWare.New(jwtWare.Config{
		ContextKey: "jwtToken",
		SigningKey: []byte("secret"),
	})
}

func MiddlewareAuthUser(db db, cache *cache.Cache) fiber.Handler {
	auth := NewAuth(db, cache)

	return func(ctx *fiber.Ctx) error {
		if token := ctx.Locals("jwtToken"); token == nil {
			return ctx.Next()
		}

		user, err := auth.authUser(ctx)
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		ctx.Locals("authUser", user)

		return ctx.Next()
	}
}

type Auth struct {
	db    db
	cache *cache.Cache
}

func NewAuth(db db, cache *cache.Cache) *Auth {
	return &Auth{db: db, cache: cache}
}

func (a *Auth) CreateToken(userID string) (string, error) {
	expireAt := time.Hour * 24 * 7
	key := utils.RandomString(40)
	err := a.cache.Set(key, userID, expireAt)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"key": key,
		"exp": time.Now().Add(expireAt).Unix(),
	}

	var t string
	t, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte("secret"))
	return t, err
}

func (a *Auth) DestroyToken(ctx *fiber.Ctx) error {
	key := a.authUserCacheKey(ctx)

	err := a.cache.Forget(key)
	return err
}

func (a *Auth) authUserCacheKey(ctx *fiber.Ctx) string {
	user := ctx.Locals("jwtToken").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["key"].(string)
}

func (a *Auth) authUser(ctx *fiber.Ctx) (interface{}, error) {
	key, err := a.cache.Get(a.authUserCacheKey(ctx))

	user, err := a.db.Find(key)
	if err != nil {
		return nil, err
	}

	return user, nil
}
