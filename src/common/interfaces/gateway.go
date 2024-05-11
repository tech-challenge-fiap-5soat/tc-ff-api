package interfaces

import (
	"net/http"

	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type IGateway interface {
	Post(req *http.Request) (*http.Response, error)
}

type PaymentGateway interface {
	GetQRCodeFromOrder(orderId string, value float64) (valueobject.QRCodeResponse, error)
}
