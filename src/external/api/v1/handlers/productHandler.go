package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type productHandler struct {
	interactor interfaces.ProductController
}

func NewProductHandler(gRouter *gin.RouterGroup, interactor interfaces.ProductController) {
	handler := &productHandler{
		interactor: interactor,
	}

	gRouter.GET("/product", handler.GetAllProductsHandler)
	gRouter.GET("/product/:category", handler.GetProductByCategoryHandler)
	gRouter.POST("/product", handler.CreateProductHandler)
	gRouter.PUT("/product/:id", handler.UpdateProductHandler)
	gRouter.DELETE("/product/:id", handler.DeleteProductHandler)
}

// Get All Products godoc
// @Summary Get all products
// @Description Get all products
// @Tags Product Routes
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Product{}
// @Router /api/v1/product [get]
func (handler *productHandler) GetAllProductsHandler(c *gin.Context) {
	actions, err := handler.interactor.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actions)
}

// Get All Products by category godoc
// @Summary Get all products by category
// @Description Get all products by category
// @Tags Product Routes
// @Param        category   path      string  true  "acompanhamento, bebida, lanche or sobremesa"
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Product{}
// @Router /api/v1/product/{category} [get]
func (handler *productHandler) GetProductByCategoryHandler(c *gin.Context) {
	supposedCategory, exists := c.Params.Get("category")
	category := valueobject.Category(supposedCategory)

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	if !valueobject.Category(category).IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	products, err := handler.interactor.GetByCategory(valueobject.Category(category))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Create Product godoc
// @Summary Create new product
// @Description Create new product
// @Tags Product Routes
// @Param        data   body      entity.Product  true  "Product information"
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/v1/product [post]
func (handler *productHandler) CreateProductHandler(c *gin.Context) {
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := valueobject.Category(product.Category)

	if !valueobject.Category(category).IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	if !product.IsValidPrice() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	if !product.IsValidName() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}

	err := handler.interactor.Create(&product)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Update Product godoc
// @Summary Update product
// @Description Update product
// @Tags Product Routes
// @Param        id   path      string  true  "Product ID"
// @Param        data   body      entity.Product  true  "Product information"
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/v1/product/{id} [put]
func (handler *productHandler) UpdateProductHandler(c *gin.Context) {
	productId, exists := c.Params.Get("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := valueobject.Category(product.Category)

	if !valueobject.Category(category).IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	if !product.IsValidPrice() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	if !product.IsValidName() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}

	err := handler.interactor.Update(productId, &product)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Delete Product godoc
// @Summary Delete product
// @Description Delete product
// @Tags Product Routes
// @Param        id   path      string  true  "Product ID"
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/v1/product/{id} [delete]
func (handler *productHandler) DeleteProductHandler(c *gin.Context) {
	productId, exists := c.Params.Get("id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	err := handler.interactor.Delete(productId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
