package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Srv struct {
	baseURL  string
	e        *echo.Echo
	redisSvc *redis.Client
}

func main() {
	srv := &Srv{
		baseURL: "http://localhost:5564",
		e:       echo.New(),
		redisSvc: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	srv.e.GET("/s/:id", srv.RedirectToLongURL)
	srv.e.POST("/v1/surl", srv.CreateShortURL)
	srv.e.Use(middleware.Logger())
	srv.e.Start(":5564")
}

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url"`
}

type CreateShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

func (srv *Srv) CreateShortURL(ctx echo.Context) error {
	res := &CreateShortURLRequest{}
	if err := ctx.Bind(res); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err)
	}

	if res.OriginalURL == "" {
		return ctx.JSON(http.StatusUnprocessableEntity, "Missing original_url field.")
	}

	randStr := generateRandomString(6)

	if ok, err := srv.redisSvc.SetNX(randStr, res.OriginalURL, time.Minute*5).Result(); !ok {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, CreateShortURLResponse{
		ShortURL: srv.baseURL + "/s/" + randStr,
	})
}

func generateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (srv *Srv) RedirectToLongURL(ctx echo.Context) error {
	val, err := srv.redisSvc.Get(ctx.Param("id")).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.Redirect(http.StatusSeeOther, val)
}
