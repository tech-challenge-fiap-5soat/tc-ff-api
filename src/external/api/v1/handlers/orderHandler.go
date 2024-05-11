package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	orderStatus "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type orderHandler struct {
	interactor interfaces.OrderController
}

func NewOrderHandler(gRouter *gin.RouterGroup, interactor interfaces.OrderController) {
	handler := &orderHandler{
		interactor: interactor,
	}

	gRouter.GET("/order", handler.FindAllHandler)
	gRouter.GET("/order/:id", handler.FindByIdHandler)
	gRouter.GET("/order/status/:status", handler.GetAllByStatusHandler)
	gRouter.POST("/order", handler.CreateOrderHandler)
	gRouter.PUT("/order/:id", handler.UpdateOrderHandler)
	gRouter.PUT("/order/:id/status/:status", handler.UpdateStatusOrderHandler)
}

// Get All Orders godoc
// @Summary Get all orders
// @Description Get all orders
// @Tags Order Routes
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Order{}
// @Router /api/v1/order [get]
func (handler *orderHandler) FindAllHandler(c *gin.Context) {
	orders, err := handler.interactor.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Get Order godoc
// @Summary Get order by ID
// @Description Get order by ID
// @Tags Order Routes
// @Param        id   path      string  true  "Order ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.Order{}
// @Router /api/v1/order/{id} [get]
func (handler *orderHandler) FindByIdHandler(c *gin.Context) {
	orderId, exists := c.Params.Get("id")
	var order *entity.Order // Only to swaggo doc

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	result, err := handler.interactor.FindById(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order = result

	c.JSON(http.StatusOK, order)
}

// Get All Orders by Status godoc
// @Summary Get all orders by status
// @Description Get all orders by status
// @Tags Order Routes
// @Param        status   path      string  true  "STARTED, PAYMENT_PENDING, PAYMENT_APPROVED, PAYMENT_REFUSED, PREPARING, READY or COMPLETED"
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Order{}
// @Router /api/v1/order/status/{status} [get]
func (handler *orderHandler) GetAllByStatusHandler(c *gin.Context) {
	status, exists := c.Params.Get("status")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	orderSts, err := orderStatus.ParseOrderStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	orders, err := handler.interactor.GetAllByStatus(orderSts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Create Order godoc
// @Summary Create new order
// @Description Create new order
// @Tags Order Routes
// @Param        data   body      dto.OrderCreateDTO  true  "Order information and customer CPF"
// @Accept  json
// @Produce  json
// @Success 200 {object} interface{}
// @Router /api/v1/order [post]
func (handler *orderHandler) CreateOrderHandler(c *gin.Context) {
	var order dto.OrderCreateDTO

	if err := c.ShouldBindJSON(&order); err != nil {
		var verr validator.ValidationErrors
		var msgFieldError string
		if errors.As(err, &verr) {
			msgFieldError = strings.Split(verr[0].Namespace(), ".")[1]
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(fieldErrMsg, msgFieldError)})
			return
		}
	}

	if len(order.OrderItemsDTO) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "at least one product must be selected"})
		return
	}

	orderId, err := handler.interactor.CreateOrder(order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"OrderId": orderId})
}

// Update Order godoc
// @Summary Update order
// @Description Update order
// @Tags Order Routes
// @Param        data   body      dto.OrderUpdateDTO  true  "Order information and customer CPF"
// @Accept  json
// @Produce  json
// @Success 204
// @Router /api/v1/order/{id} [put]
func (handler *orderHandler) UpdateOrderHandler(c *gin.Context) {
	orderId, exists := c.Params.Get("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	var order dto.OrderUpdateDTO
	if err := c.ShouldBindJSON(&order); err != nil {
		var verr validator.ValidationErrors
		var msgFieldError string
		if errors.As(err, &verr) {
			msgFieldError = strings.Split(verr[0].Namespace(), ".")[1]
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(fieldErrMsg, msgFieldError)})
			return
		}
	}

	if len(order.OrderItemsDTO) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "at least one product must be selected"})
		return
	}

	err := handler.interactor.UpdateOrder(orderId, order)

	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// Update Order Status godoc
// @Summary Update order status
// @Description Update order status
// @Tags Order Routes
// @Param        id   path      string  true  "Order ID"
// @Param        status   path      string  true  "STARTED, PREPARING, READY or COMPLETED"
// @Accept  json
// @Produce  json
// @Success 204
// @Router /api/v1/order/{id}/status/{status} [put]
func (handler orderHandler) UpdateStatusOrderHandler(c *gin.Context) {
	orderId, exists := c.Params.Get("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	status, exists := c.Params.Get("status")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	sts, err := orderStatus.ParseOrderStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	if sts.IsPaymentStatus() {
		c.JSON(http.StatusForbidden, gin.H{"error": "it is not allowed to directly change the payment status"})
		return
	}

	err = handler.interactor.UpdateOrderStatus(orderId, sts)

	if err != nil && strings.Contains(err.Error(), "record not found") {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if err != nil && strings.Contains(err.Error(), "cannot be updated to") {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if err != nil && strings.Contains(err.Error(), "previous status") {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
