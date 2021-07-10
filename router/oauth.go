package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal/oauth"
)

// AuthResponse 認証の返答

func (h *Handlers) CallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return c.String(http.StatusBadRequest, "Code Is Null")
	}
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err).Error()) //TODO:エラーを返さないように
	}
	codeVerifier, ok := sess.Values["codeVerifier"].(string)
	if !ok {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err).Error())
	}
	res, err := h.authConf.GetToken(code, codeVerifier)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get access token: %w", err).Error())
	}
	sess.Values["accsessToken"] = res.AccessToken
	sess.Values["refreshToken"] = res.RefreshToken
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   res.Expiry.Second(),
		HttpOnly: true,
	}

	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func (h *Handlers) PostGenerateCodeHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err).Error())
	}
	pkce, err := oauth.GenerateCode()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get access token:%w", err).Error())
	}
	sess.Values["codeVerifier"] = pkce.CodeVerifier
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 1000,
		HttpOnly: true,
	}
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to set access token:%w", err).Error())
	}

	return c.JSON(http.StatusOK, pkce)
}
