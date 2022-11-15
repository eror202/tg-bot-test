package handler

import (
	"github.com/gin-gonic/gin"
	"tg-bot-for-ts/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	botMsg := router.Group("/messages")
	{
		botMsg.GET("/msg", h.getAllMessage)
		botMsg.POST("/msg", h.createMessage)
	}
	return router
}
