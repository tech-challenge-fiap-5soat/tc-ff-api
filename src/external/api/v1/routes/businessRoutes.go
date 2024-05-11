package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/constants"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	"github.com/hcsouza/fiap-tech-fast-food/src/external/api/infra/config"
	"github.com/hcsouza/fiap-tech-fast-food/src/external/api/v1/handlers"
	"github.com/hcsouza/fiap-tech-fast-food/src/operation/controller"
	"github.com/hcsouza/fiap-tech-fast-food/src/operation/gateway"

	mongodb "github.com/hcsouza/fiap-tech-fast-food/src/external/datasource"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterBusinessRoutes(gServer *gin.RouterGroup, dbClient mongo.Client) {
	groupServer := gServer.Group("/v1")

	registerCustomerHandler(groupServer, dbClient)
	registerProductHandler(groupServer, dbClient)
	registerOrderHandler(groupServer, dbClient)
	registerCheckoutHandler(groupServer, dbClient)
}

func registerCustomerHandler(groupServer *gin.RouterGroup, dbClient mongo.Client) {
	mongoAdapter := mongodb.NewMongoAdapter[entity.Customer](
		dbClient,
		config.GetMongoCfg().Database,
		constants.CustomerCollection,
	)

	customerInteractor := controller.NewCustomerController(mongoAdapter)
	handlers.NewCustomerHandler(groupServer, customerInteractor)
}

func registerProductHandler(groupServer *gin.RouterGroup, dbClient mongo.Client) {
	mongoAdapter := mongodb.NewMongoAdapter[entity.Product](
		dbClient,
		config.GetMongoCfg().Database,
		constants.ProductCollection,
	)

	productInteractor := controller.NewProductController(mongoAdapter)
	handlers.NewProductHandler(groupServer, productInteractor)
}

func registerOrderHandler(groupServer *gin.RouterGroup, dbClient mongo.Client) {
	orderDbAdapter := mongodb.NewMongoAdapter[entity.Order](
		dbClient,
		config.GetMongoCfg().Database,
		constants.OrderCollection,
	)

	productDbAdapter := mongodb.NewMongoAdapter[entity.Product](
		dbClient,
		config.GetMongoCfg().Database,
		constants.ProductCollection,
	)
	productGateway := gateway.NewProductGateway(productDbAdapter)
	productUseCase := usecase.NewProductUseCase(productGateway)

	customerDbAdapter := mongodb.NewMongoAdapter[entity.Customer](
		dbClient,
		config.GetMongoCfg().Database,
		constants.CustomerCollection,
	)
	customerGateway := gateway.NewCustomerGateway(customerDbAdapter)
	customerUseCase := usecase.NewCustomerUseCase(customerGateway)

	orderInteractor := controller.NewOrderController(orderDbAdapter, productUseCase, customerUseCase)

	handlers.NewOrderHandler(groupServer, orderInteractor)
}

func registerCheckoutHandler(groupServer *gin.RouterGroup, dbClient mongo.Client) {

	productDbAdapter := mongodb.NewMongoAdapter[entity.Product](
		dbClient,
		config.GetMongoCfg().Database,
		constants.ProductCollection,
	)
	productGateway := gateway.NewProductGateway(productDbAdapter)
	productUseCase := usecase.NewProductUseCase(productGateway)

	customerDbAdapter := mongodb.NewMongoAdapter[entity.Customer](
		dbClient,
		config.GetMongoCfg().Database,
		constants.CustomerCollection,
	)
	customerGateway := gateway.NewCustomerGateway(customerDbAdapter)
	customerUseCase := usecase.NewCustomerUseCase(customerGateway)

	orderDbAdapter := mongodb.NewMongoAdapter[entity.Order](
		dbClient,
		config.GetMongoCfg().Database,
		constants.OrderCollection,
	)
	orderGateway := gateway.NewOrderGateway(orderDbAdapter)
	orderUseCase := usecase.NewOrderUseCase(orderGateway, productUseCase, customerUseCase)

	checkoutInteractor := controller.NewCheckoutController(orderUseCase)

	handlers.NewCheckoutHandler(groupServer, checkoutInteractor)
}
