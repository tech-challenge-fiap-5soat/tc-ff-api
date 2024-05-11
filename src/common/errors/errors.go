package errors

import "errors"

var ErrDuplicatedKey = errors.New("duplicate key error on collection")
var ErrCheckoutOrderAlreadyCompleted = errors.New("order already has a checkout completed")
var ErrInvalidCategory = errors.New("invalid category")
