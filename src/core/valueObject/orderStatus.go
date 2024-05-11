package valueobject

import (
	"fmt"
	"slices"
	"strings"
)

type OrderStatus string

const (
	ORDER_STARTED          OrderStatus = "STARTED"
	ORDER_PAYMENT_PENDING  OrderStatus = "PAYMENT_PENDING"
	ORDER_PAYMENT_APPROVED OrderStatus = "PAYMENT_APPROVED"
	ORDER_PAYMENT_REFUSED  OrderStatus = "PAYMENT_REFUSED"
	ORDER_BEING_PREPARED   OrderStatus = "PREPARING"
	ORDER_READY            OrderStatus = "READY"
	ORDER_COMPLETED        OrderStatus = "COMPLETED"
)

func (o OrderStatus) String() string {
	return string(o)
}

func ParseOrderStatus(s string) (o OrderStatus, err error) {
	statuses := map[OrderStatus]struct{}{
		ORDER_STARTED:          {},
		ORDER_PAYMENT_PENDING:  {},
		ORDER_PAYMENT_APPROVED: {},
		ORDER_PAYMENT_REFUSED:  {},
		ORDER_BEING_PREPARED:   {},
		ORDER_READY:            {},
		ORDER_COMPLETED:        {},
	}

	orderStatus := OrderStatus(strings.ToUpper(s))
	_, ok := statuses[orderStatus]

	if !ok {
		return o, fmt.Errorf(`cannot parse:[%s] as order status`, s)
	}
	return orderStatus, nil
}

func (o OrderStatus) AvailableNextStatus(status OrderStatus) []OrderStatus {
	switch status {
	case ORDER_STARTED:
		return []OrderStatus{ORDER_PAYMENT_PENDING}
	case ORDER_PAYMENT_PENDING:
		return []OrderStatus{ORDER_PAYMENT_APPROVED, ORDER_PAYMENT_REFUSED}
	case ORDER_PAYMENT_APPROVED:
		return []OrderStatus{ORDER_BEING_PREPARED}
	case ORDER_BEING_PREPARED:
		return []OrderStatus{ORDER_READY}
	case ORDER_READY:
		return []OrderStatus{ORDER_COMPLETED}
	default:
		return []OrderStatus{}
	}
}

func (o OrderStatus) IsValidNextStatus(nextStatus string) bool {
	currentStatus, err := ParseOrderStatus(o.String())
	allowSameStatus := false

	if err != nil {
		return false
	}

	nextStatusParsed, err := ParseOrderStatus(nextStatus)

	if err != nil {
		return false
	}

	if nextStatusParsed == currentStatus {
		return allowSameStatus
	}

	return slices.Contains(currentStatus.AvailableNextStatus(currentStatus), nextStatusParsed)
}

func (o OrderStatus) IsPaymentStatus() bool {
	return o == ORDER_PAYMENT_PENDING || o == ORDER_PAYMENT_APPROVED || o == ORDER_PAYMENT_REFUSED
}

func (o OrderStatus) GetPreviousStatus() []OrderStatus {
	switch o {
	case ORDER_PAYMENT_PENDING:
		return []OrderStatus{ORDER_STARTED}
	case ORDER_PAYMENT_APPROVED:
		return []OrderStatus{ORDER_STARTED, ORDER_PAYMENT_PENDING}
	case ORDER_PAYMENT_REFUSED:
		return []OrderStatus{ORDER_STARTED, ORDER_PAYMENT_PENDING}
	case ORDER_BEING_PREPARED:
		return []OrderStatus{ORDER_STARTED, ORDER_PAYMENT_PENDING, ORDER_PAYMENT_APPROVED, ORDER_PAYMENT_REFUSED}
	case ORDER_READY:
		return []OrderStatus{ORDER_STARTED, ORDER_PAYMENT_PENDING, ORDER_PAYMENT_APPROVED, ORDER_PAYMENT_REFUSED, ORDER_BEING_PREPARED}
	case ORDER_COMPLETED:
		return []OrderStatus{ORDER_STARTED, ORDER_PAYMENT_PENDING, ORDER_PAYMENT_APPROVED, ORDER_PAYMENT_REFUSED, ORDER_BEING_PREPARED, ORDER_READY}
	default:
		return []OrderStatus{}
	}
}

func (o OrderStatus) OrderCanBeUpdated() bool {
	return o == ORDER_STARTED
}
