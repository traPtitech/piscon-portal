package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func middlewareAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("failed to get session:%w", err))
		}

		accessToken := sess.Values["accessToken"]
		if accessToken == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set("accessToken", accessToken)

		return next(c)
	}
}
