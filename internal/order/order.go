package order

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("order not found")

type Order struct {
	ID          string    `json:"id"`
	AmountCents int64     `json:"amount_cents"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store interface {
	Get(ctx context.Context, id string) (Order, error)
	Create(ctx context.Context, o Order) error
}

type MemStore struct {
	mu     sync.Mutex
	orders map[string]Order
}

func NewMemStore() *MemStore {
	return &MemStore{
		orders: make(map[string]Order),
	}
}

func (s *MemStore) Get(ctx context.Context, id string) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]

	if !ok {
		return Order{}, ErrNotFound
	}

	return o, nil
}

func (s *MemStore) Create(ctx context.Context, o Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[o.ID] = o

	return nil
}
