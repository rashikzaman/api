package http

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func RunHTTPServer(app application.Application, redisStore redis.Store) {
	r := gin.Default()
	r.Use(sessions.Sessions("user_session", redisStore))

	routes.AuthRoutes(r, app)

	_ = r.Run(":" + app.Config.GetHTTPPort())
}
