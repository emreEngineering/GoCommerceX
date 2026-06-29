package adapters

import (
	"context"

	"GoCommerceX/order-service/internal/ports"
	paymentv1 "GoCommerceX/proto/payment/v1"
)

type PaymentServiceClient struct {
	client paymentv1.PaymentServiceClient
}

func NewPaymentServiceClient(client paymentv1.PaymentServiceClient) *PaymentServiceClient {
	return &PaymentServiceClient{client: client}
}

func (c *PaymentServiceClient) CreatePayment(ctx context.Context, req ports.PaymentRequest) (ports.PaymentResult, error) {
	resp, err := c.client.CreatePayment(ctx, &paymentv1.CreatePaymentRequest{
		OrderId: req.OrderID,
		Amount:  req.Amount,
		Method:  req.Method,
	})
	if err != nil {
		return ports.PaymentResult{}, err
	}
	return ports.PaymentResult{
		PaymentID: resp.GetPayment().GetId(),
		Status:    resp.GetPayment().GetStatus(),
	}, nil
}
