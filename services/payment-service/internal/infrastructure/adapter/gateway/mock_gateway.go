// Mock payment gateway adapter (simulates Stripe)
package gateway

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/cloudmart/payment-service/internal/domain/model"
	"github.com/cloudmart/payment-service/internal/domain/port"
)

type mockGateway struct{}

func NewMockPaymentGateway() port.PaymentGateway {
	return &mockGateway{}
}

func (g *mockGateway) Charge(ctx context.Context, amount float64, currency string, method model.PaymentMethod, metadata map[string]interface{}) (string, error) {
	// Simulate processing delay
	time.Sleep(100 * time.Millisecond)

	// Simulate 95% success rate
	if rand.Float64() < 0.05 {
		return "", fmt.Errorf("payment declined by issuer")
	}

	txID := fmt.Sprintf("pi_%d_%d", time.Now().UnixNano(), rand.Intn(99999))
	return txID, nil
}

func (g *mockGateway) Refund(ctx context.Context, providerTxID string) error {
	time.Sleep(50 * time.Millisecond)
	return nil
}
