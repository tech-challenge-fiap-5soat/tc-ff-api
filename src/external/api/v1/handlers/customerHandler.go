package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	coreErrors "github.com/hcsouza/fiap-tech-fast-food/src/common/errors"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	vo "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
	"github.com/hcsouza/fiap-tech-fast-food/src/external/auth"
)

type customerHandler struct {
	interactor interfaces.CustomerController
}

const (
	fieldErrMsg = "Invalid value for field: '%s'"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("IsCpfValid", CpfValidator)
	}
}

func NewCustomerHandler(gRouter *gin.RouterGroup, interactor interfaces.CustomerController) {
	handler := &customerHandler{
		interactor: interactor,
	}

	gRouter.GET("/customer", handler.GetCustomerHandler)
	gRouter.POST("/customer", handler.CreateCustomerHandler)
	gRouter.GET("customer/authorization", handler.GetAuthorizationTokenHandler)

}

// Create Customer godoc
// @Summary Create a new customer
// @Description Create a new customer
// @Tags Customer Routes
// @Param        data    body     dto.CustomerCreateDTO  true  "Customer information"
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.Customer{}
// @Router /api/v1/customer [post]
func (handler *customerHandler) CreateCustomerHandler(ctx *gin.Context) {
	var createRequest dto.CustomerCreateDTO
	err := ctx.ShouldBindJSON(&createRequest)

	if err != nil {
		var verr validator.ValidationErrors
		var msgFieldError string
		if errors.As(err, &verr) {
			msgFieldError = strings.Split(verr[0].Namespace(), ".")[1]
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(fieldErrMsg, msgFieldError)})
			return
		}
	}

	customer, err := handler.interactor.CreateCustomer(ctx.Request.Context(), createRequest)

	if errors.Is(err, coreErrors.ErrDuplicatedKey) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "customer already exists"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, customer)
}

// Get Customer godoc
// @Summary Get customer by CPF
// @Description Get customer by CPF
// @Tags Customer Routes
// @Param        cpf    query     string  true  "19119119100"
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Customer{}
// @Router /api/v1/customer [get]
func (handler *customerHandler) GetCustomerHandler(ctx *gin.Context) {
	cpf := ctx.Query("cpf")
	params := map[string]string{"cpf": cpf}
	var customer *entity.Customer // Only to swaggo doc

	actions, err := handler.interactor.GetCustomer(ctx.Request.Context(), params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customer = actions

	ctx.JSON(http.StatusOK, customer)
}

// Get AuthToken godoc
// @Summary Get authorization token by CPF
// @Description Get authorization tokenCPF
// @Tags Customer Routes
// @Param        cpf    query     string  true  "19119119100"
// @Accept  json
// @Produce  json
// @Success 200 {array} auth.AuthorizationToken{}
// @Router /api/v1/customer/authorization [get]
func (handler *customerHandler) GetAuthorizationTokenHandler(ctx *gin.Context) {
	cpf := ctx.Query("cpf")
	var authToken *auth.AuthorizationToken // Only to swaggo doc

	authToken, err := auth.GetAuthorizationToken(cpf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(authToken.AccessToken) < 1 || err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate authorization token"})
		return
	}

	ctx.JSON(http.StatusOK, authToken)
}

func CpfValidator(fl validator.FieldLevel) bool {
	rawCpf, ok := fl.Field().Interface().(string)
	cpfToValidate := vo.CPF(rawCpf)

	if ok {
		return cpfToValidate.IsValid()
	}
	return false
}
