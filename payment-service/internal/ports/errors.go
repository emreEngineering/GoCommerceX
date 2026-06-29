package ports

import "errors"

var ErrPaymentNotFound = errors.New("payment repository: payment not found")
