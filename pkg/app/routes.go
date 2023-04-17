package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	r := s.router

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	if s.expressionHandler != nil {
		expGroup := r.Group("/expressions")
		expGroup.POST("/", s.expressionHandler.CreateExpression)
		expGroup.GET("/", s.expressionHandler.ListExpressions)
		expGroup.PUT("/:id", s.expressionHandler.UpdateExpression)

		r.GET("/evaluate/:id", s.expressionHandler.EvaluateExpression)
	}

	return r
}
