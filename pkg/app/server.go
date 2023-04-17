package app

import (
	"fmt"
	"log"

	"github.com/CaioTeixeira95/logic-exp/pkg/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router            *gin.Engine
	expressionHandler *handlers.ExpressionHandler
}

type ServerHandlerOption func(s *Server)

func NewServer(router *gin.Engine, handlerOptions ...ServerHandlerOption) *Server {
	s := &Server{
		router: router,
	}

	for _, handlerOption := range handlerOptions {
		handlerOption(s)
	}

	return s
}

func (s *Server) Run(addr ...string) error {
	r := s.Routes()
	if err := r.Run(addr...); err != nil {
		log.Printf("Server - there was an error calling Run on router: %s", err.Error())
		return fmt.Errorf("error trying to Run server: %w", err)
	}
	return nil
}

func WithExpressionHandlerOption(h *handlers.ExpressionHandler) ServerHandlerOption {
	return func(s *Server) {
		s.expressionHandler = h
	}
}
