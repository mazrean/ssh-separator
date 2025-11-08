package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mazrean/separated-webshell/api/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

var (
	prometheus = os.Getenv("PROMETHEUS")
)

type API struct {
	*User
}

func NewAPI(user *User) *API {
	return &API{
		User: user,
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

	if prometheus == "true" {
		e.Use(middlewares.RequestCounter())
		e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	}

	// API認証エンドポイント用のレート制限設定
	// 環境変数から設定を読み取る（デフォルト: 1分間に5回まで）
	rateLimitRate := 5
	if envRate, ok := os.LookupEnv("RATE_LIMIT_RATE"); ok {
		if parsedRate, err := strconv.Atoi(envRate); err == nil && parsedRate > 0 {
			rateLimitRate = parsedRate
		}
	}

	rateLimitBurst := 5
	if envBurst, ok := os.LookupEnv("RATE_LIMIT_BURST"); ok {
		if parsedBurst, err := strconv.Atoi(envBurst); err == nil && parsedBurst > 0 {
			rateLimitBurst = parsedBurst
		}
	}

	rateLimitExpiresIn := int(time.Minute / time.Second)
	if envExpires, ok := os.LookupEnv("RATE_LIMIT_EXPIRES_IN"); ok {
		if parsedExpires, err := strconv.Atoi(envExpires); err == nil && parsedExpires > 0 {
			rateLimitExpiresIn = parsedExpires
		}
	}

	rateLimiterConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(rateLimitRate),
				Burst:     rateLimitBurst,
				ExpiresIn: time.Duration(rateLimitExpiresIn) * time.Second,
			},
		),
	}

	e.POST("/new", api.PostNewUser, middleware.RateLimiterWithConfig(rateLimiterConfig))
	e.PUT("/reset", api.PutReset, middleware.RateLimiterWithConfig(rateLimiterConfig))

	return e.Start(fmt.Sprintf(":%d", port))
}
