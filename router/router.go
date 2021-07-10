package router

import (
	"github.com/labstack/echo/v4"
)

func (h *Handlers) SetUp(e *echo.Echo) {
	api := e.Group("/api")
	api.POST("/auth/code", h.PostGenerateCodeHandler)
	api.GET("/auth/callback", h.CallbackHandler)
	api.GET("/results", h.GetAllResults)
	api.GET("/benchmark/queue", h.GetBenchmarkQueue)
	api.GET("/newer", h.GetNewer)
	api.GET("/questions", h.GetQuestions)
	// api.POST("/instancelog", postInstanceLog)

	apiWithAuth := e.Group("/api", middlewareAuthUser)
	apiWithAuth.GET("/me", h.GetMeFromTraq)
	apiWithAuth.GET("/me/group", h.GetMeGroup)
	apiWithAuth.POST("/team", h.CreateTeam)
	apiWithAuth.POST("/user", h.CreateUser)
	apiWithAuth.POST("/instance/:team_id/:instance_number", h.CreateInstance)
	apiWithAuth.DELETE("/instance/:team_id/:instance_number", h.DeleteInstance)
	// TODO: ユーザー名で認証してないので修正する必要がある
	apiWithAuth.GET("/team/:id", h.GetTeam)
	apiWithAuth.GET("/user/:name", h.GetUser)
	apiWithAuth.POST("/benchmark/:name/:instance_number", h.QueBenchmark)
	apiWithAuth.GET("/admin/team", h.GetAllTeam)

	apiWithAuth.POST("/questions", h.PostQuestions)
	apiWithAuth.PUT("/questions/:id", h.PutQuestions)
	apiWithAuth.DELETE("/questions/:id", h.DeleteQuestions)

	// e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(os.Getenv("HOST"))
	// e.AutoTLSManager.Cache = autocert.DirCache("/etc/letsencrypt/live/piscon-portal.trap.jp/cert.pem")
	// e.Pre(middleware.HTTPSWWWRedirect())
	// switch env {
	// case "prod":
	// e.StartAutoTLS(":443")
	// default:
}
