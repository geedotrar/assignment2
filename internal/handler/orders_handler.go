package handler

import (
	"net/http"
	"strconv"

	"assignment2/internal/model"
	"assignment2/internal/service"
	"assignment2/pkg"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	GetOrders(ctx *gin.Context)
	GetOrdersById(ctx *gin.Context)
	CreateOrder(ctx *gin.Context)
	DeleteOrder(ctx *gin.Context)
	UpdateOrder(ctx *gin.Context)
}

type orderHandlerImpl struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) OrderHandler {
	return &orderHandlerImpl{
		svc: svc,
	}
}

func (u *orderHandlerImpl) GetOrders(ctx *gin.Context) {

	orders, err := u.svc.GetOrders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

func (u *orderHandlerImpl) GetOrdersById(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	order, err := u.svc.GetOrderById(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (u *orderHandlerImpl) CreateOrder(ctx *gin.Context) {
	var requestData model.Order

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
		return
	}

	createdOrder, err := u.svc.CreateOrder(ctx, requestData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdOrder)
}

func (u *orderHandlerImpl) DeleteOrder(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid order ID"})
		return
	}

	err = u.svc.DeleteOrder(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func (u *orderHandlerImpl) UpdateOrder(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid order ID"})
		return
	}

	var updatedOrder model.Order
	if err := ctx.ShouldBindJSON(&updatedOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
		return
	}

	err = u.svc.UpdateOrder(ctx, id, updatedOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}
