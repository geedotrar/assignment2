package router

import (
	"assignment2/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Mount()
}

type orderRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.OrderHandler
}

func NewOrderRouter(v *gin.RouterGroup, handler handler.OrderHandler) UserRouter {
	return &orderRouterImpl{v: v, handler: handler}
}

func (u *orderRouterImpl) Mount() {
	u.v.GET("", u.handler.GetOrders)
	u.v.GET("/:id", u.handler.GetOrdersById)
	u.v.POST("", u.handler.CreateOrder)
	u.v.DELETE("/:id", u.handler.DeleteOrder)
	u.v.PUT("/:id", u.handler.UpdateOrder)
}
