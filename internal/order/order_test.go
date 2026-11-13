package order_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/yashcherba/pet-api/internal/order"
)

func TestMemStoreGet(t *testing.T) {
	tests := []struct {
		name    string
		seed    bool
		wantErr error
	}{
		{
			name:    "found",
			seed:    true,
			wantErr: nil,
		},
		{
			name:    "not found",
			seed:    false,
			wantErr: order.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			st := order.NewMemStore()

			const id = "order-1"
			if tt.seed {
				if err := st.Create(context.Background(), order.Order{ID: id}); err != nil {
					t.Fatalf("seed Create: %v", err)
				}
			}

			_, err := st.Get(context.Background(), id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStoreConcurrent(t *testing.T) {
	st := order.NewMemStore()

	const n = 100

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			id := fmt.Sprintf("order-%d", i%10)

			if i%2 == 0 {
				_ = st.Create(context.Background(), order.Order{ID: id})
			} else {
				_, _ = st.Get(context.Background(), id)
			}
		}()
	}

	wg.Wait()
}
