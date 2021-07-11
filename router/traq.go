package router

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal/model"
)

var baseURL, _ = url.Parse("https://q.trap.jp/api/v3")

func (h *Handlers) GetMeFromTraq(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server error")
	}
	token := sess.Values["accessToken"].(string)
	u, err := getMe(token)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server error")
	}

	return c.JSON(http.StatusOK, model.User{
		Name:       u.Id,
		ScreenName: u.Name,
		IconFileId: u.Icon,
	})

}

func (h *Handlers) GetMeGroup(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server error")
	}
	token := sess.Values["accessToken"].(string)
	u, err := getMe(token)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server error")
	}
	path := *baseURL
	gid := u.Groups[0]
	path.Path = baseURL.Path + "/groups/" + gid
	req, err := http.NewRequest("GET", path.String(), nil)

	req.Header.Set("Authorization", "Bearer "+token)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server error")
	}
	if res.StatusCode != http.StatusOK {
		return c.String(res.StatusCode, "Failed to send Request")
	}
	var g struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(res.Body).Decode(&g)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server error")
	}
	return c.String(http.StatusOK, g.Name)

}

type myUserDetail struct {
	Id     string   `json:"name"`
	Name   string   `json:"displayName"`
	Icon   string   `json:"iconFileId"`
	Groups []string `json:"groups"`
}

func getMe(token string) (*myUserDetail, error) {
	path := *baseURL
	path.Path += "/users/me"
	req, err := http.NewRequest("GET", path.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, err
	}
	var u myUserDetail
	err = json.NewDecoder(res.Body).Decode(&u)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
