package facade

import (
	"fmt"
)

type Order struct {
	ID     string
	UserID string
	ItemID string
	Amount int64
}

// Системы

type PaymentService struct{}

func (p PaymentService) Charge(userID string, amount int64) error {
	fmt.Printf("Charge user %s for %d\n", userID, amount)

	return nil
}

type InventoryService struct{}

func (i InventoryService) Reserve(itemID string) error {
	fmt.Printf("Reserve item %s on stock\n", itemID)

	return nil
}

type FraudService struct{}

func (f FraudService) Check(order Order) error {
	fmt.Printf("Run fraud check for order %s\n", order.ID)

	return nil
}

// Фасад

type OrderService struct {
	payment   PaymentService
	inventory InventoryService
	fraud     FraudService
}

func NewOrderService() *OrderService {
	return &OrderService{
		payment:   PaymentService{},
		inventory: InventoryService{},
		fraud:     FraudService{},
	}
}

func (s *OrderService) PlaceOrder(o Order) error {
	if err := s.fraud.Check(o); err != nil {
		return fmt.Errorf("fraud check failed: %w", err)
	}

	if err := s.payment.Charge(o.UserID, o.Amount); err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}

	if err := s.inventory.Reserve(o.ItemID); err != nil {
		return fmt.Errorf("inventory failed: %w", err)
	}

	fmt.Printf("Order %s successfully placed\n", o.ID)
	return nil
}

func Example() {
	service := NewOrderService()

	order := Order{
		ID:     "order-123",
		UserID: "user-42",
		ItemID: "item-777",
		Amount: 1999,
	}

	if err := service.PlaceOrder(order); err != nil {
		fmt.Println("order error:", err)
		return
	}
}
