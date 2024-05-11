package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	coreErrors "github.com/hcsouza/fiap-tech-fast-food/src/common/errors"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	vo "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type checkoutHandler struct {
	interactor interfaces.CheckoutController
}

func NewCheckoutHandler(gRouter *gin.RouterGroup, interactor interfaces.CheckoutController) {
	handler := &checkoutHandler{
		interactor: interactor,
	}

	gRouter.POST("/checkout/:id", handler.CreateCheckout)
	gRouter.POST("/checkout/:id/callback", handler.UpdateCheckoutCallback)
}

// Create Checkout From Order godoc
// @Summary Create checkout from order
// @Description Create checkout from order
// @Tags Checkout Routes
// @Param        id   path      string  true  "Order ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.CreateCheckout{}
// @Router /api/v1/checkout/:id [post]
func (handler checkoutHandler) CreateCheckout(c *gin.Context) {
	var response dto.CreateCheckout

	orderId, exists := c.Params.Get("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	createdCheckout, err := handler.interactor.CreateCheckout(orderId)

	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response = *createdCheckout

	c.JSON(http.StatusOK, response)
}

// Update Checkout Callback godoc
// @Summary Update checkout callback
// @Description Update checkout callback
// @Tags Checkout Routes
// @Param        id   path      string  true  "Order ID"
// @Param        data   body      dto.UpdateCheckoutDTO  true  "Order payment result status: approved, refused."
// @Accept  json
// @Produce  json
// @Success 204 {object} interface{}
// @Router /api/v1/checkout/:id/callback [post]
func (handler checkoutHandler) UpdateCheckoutCallback(c *gin.Context) {
	var body dto.UpdateCheckoutDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderId, exists := c.Params.Get("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	parseOrderStatus, err := vo.ParseOrderStatus(fmt.Sprintf("PAYMENT_%s", strings.ToUpper(body.Status)))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = handler.interactor.UpdateCheckout(orderId, parseOrderStatus)

	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if err != nil && errors.Is(err, coreErrors.ErrCheckoutOrderAlreadyCompleted) {
		c.JSON(http.StatusConflict, gin.H{"error": coreErrors.ErrCheckoutOrderAlreadyCompleted.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
