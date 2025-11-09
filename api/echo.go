package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mazrean/separated-webshell/api/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

type Config struct {
	Prometheus         bool
	RateLimitRate      int
	RateLimitBurst     int
	RateLimitExpiresIn int
}

type API struct {
	*User
	config Config
}

func NewAPI(user *User, config Config) *API {
	return &API{
		User:   user,
		config: config,
	}
}

func (api *API) Start(port int) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		// Issue #1426
		code := he.Code
		message := he.Message
		if m, ok := he.Message.(string); ok {
			if e.Debug {
				message = echo.Map{"message": m, "error": err.Error()}
			} else {
				message = echo.Map{"message": m}
			}
		} else if err, ok := he.Message.(error); ok {
			c.Logger().Error(err)
			message = echo.Map{"message": err.Error()}
		}

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead { // Issue #608
				err = c.NoContent(he.Code)
			} else {
				err = c.JSON(code, message)
			}
			if err != nil {
				e.Logger.Error(err)
			}
		}
	}

	if api.config.Prometheus {
		e.Use(middlewares.RequestCounter())
		e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	}

	rateLimiterConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(api.config.RateLimitRate),
				Burst:     api.config.RateLimitBurst,
				ExpiresIn: time.Duration(api.config.RateLimitExpiresIn) * time.Second,
			},
		),
	}

	e.POST("/new", api.PostNewUser, middleware.RateLimiterWithConfig(rateLimiterConfig))
	e.PUT("/reset", api.PutReset, middleware.RateLimiterWithConfig(rateLimiterConfig))

	return e.Start(fmt.Sprintf(":%d", port))
}
