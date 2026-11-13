package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yashcherba/pet-api/internal/order"
)

func TestRouterGetOrder_Found(t *testing.T) {
	t.Parallel()

	st := order.NewMemStore()
	const id = "order-1"
	want := order.Order{ID: id, Status: "new", AmountCents: 1234}
	if err := st.Create(context.Background(), want); err != nil {
		t.Fatalf("seed Create: %v", err)
	}

	ts := httptest.NewServer(NewRouter(st))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/v1/orders/" + id)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want prefix %q", ct, "application/json")
	}

	var got order.Order
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	if got != want {
		t.Fatalf("order = %+v, want %+v", got, want)
	}
}

func TestRouterGetOrder_NotFound(t *testing.T) {
	t.Parallel()

	st := order.NewMemStore()
	ts := httptest.NewServer(NewRouter(st))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/v1/orders/does-not-exist")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestRouterHealthz_MethodNotAllowed(t *testing.T) {
	t.Parallel()

	st := order.NewMemStore()
	ts := httptest.NewServer(NewRouter(st))
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/healthz", "application/json", nil)
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	// ServeMux неявно добавляет HEAD к GET-маршрутам → Allow = "GET, HEAD".
	if allow := resp.Header.Get("Allow"); allow != "GET, HEAD" {
		t.Fatalf("Allow = %q, want %q", allow, "GET, HEAD")
	}
}
