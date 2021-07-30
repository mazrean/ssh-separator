package middlewares

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Help:      "Echo request counter.",
	Namespace: "webshell",
	Name:      "http_requests_total",
}, []string{"code", "method"})

// RequestCounter prometheus metrics用リクエストカウンター
func RequestCounter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			err = next(c)
			if c.Path() != "/metrics" {
				requestCounter.WithLabelValues(strconv.Itoa(c.Response().Status), c.Request().Method).Inc()
			}

			return err
		}
	}
}
