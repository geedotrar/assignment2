package main

import (
	"assignment2/internal/handler"
	"assignment2/internal/infrastructure"
	"assignment2/internal/repository"
	"assignment2/internal/router"
	"assignment2/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	usersGroup := g.Group("/orders")
	gorm := infrastructure.NewGormPostgres()
	orderRepo := repository.NewOrderQuery(gorm)
	orderSvc := service.NewOrderService(orderRepo)
	orderHdl := handler.NewOrderHandler(orderSvc)
	orderRouter := router.NewOrderRouter(usersGroup, orderHdl)

	orderRouter.Mount()

	g.Run(":3000")
}
