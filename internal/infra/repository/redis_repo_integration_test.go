package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func newTestRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{Addr: addr})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		t.Fatalf("failed to connect to redis at %s: %v", addr, err)
	}

	return client
}

func TestRedisRepo_SaveAndGetLastFive_TrimsAndOrders(t *testing.T) {
	client := newTestRedisClient(t)
	defer func() { _ = client.Close() }()

	repo := NewRedisRepo(client)
	userID := "user-save-trim"

	key := "passwords:" + userID
	if err := client.Del(context.Background(), key).Err(); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	for i := 1; i <= 7; i++ {
		pwd := fmt.Sprintf("pwd-%d", i)
		if err := repo.Save(userID, pwd); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	got, err := repo.GetLastFive(userID)
	if err != nil {
		t.Fatalf("GetLastFive failed: %v", err)
	}

	if len(got) != 5 {
		t.Fatalf("expected 5 passwords, got %d: %v", len(got), got)
	}

	want := []string{"pwd-7", "pwd-6", "pwd-5", "pwd-4", "pwd-3"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("at index %d expected %q, got %q (full: %v)", i, want[i], got[i], got)
		}
	}
}

func TestRedisRepo_GetLastFive_Empty(t *testing.T) {
	client := newTestRedisClient(t)
	defer func() { _ = client.Close() }()

	repo := NewRedisRepo(client)
	userID := "user-empty"

	key := "passwords:" + userID
	if err := client.Del(context.Background(), key).Err(); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	got, err := repo.GetLastFive(userID)
	if err != nil {
		t.Fatalf("GetLastFive failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d: %v", len(got), got)
	}
}

func TestRedisRepo_UserIsolation(t *testing.T) {
	client := newTestRedisClient(t)
	defer func() { _ = client.Close() }()

	repo := NewRedisRepo(client)

	userA := "user-a"
	userB := "user-b"

	keyA := "passwords:" + userA
	keyB := "passwords:" + userB
	if err := client.Del(context.Background(), keyA, keyB).Err(); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if err := repo.Save(userA, "a-1"); err != nil {
		t.Fatalf("Save(userA) failed: %v", err)
	}
	if err := repo.Save(userB, "b-1"); err != nil {
		t.Fatalf("Save(userB) failed: %v", err)
	}

	gotA, err := repo.GetLastFive(userA)
	if err != nil {
		t.Fatalf("GetLastFive(userA) failed: %v", err)
	}
	if len(gotA) != 1 || gotA[0] != "a-1" {
		t.Fatalf("unexpected userA passwords: %v", gotA)
	}

	gotB, err := repo.GetLastFive(userB)
	if err != nil {
		t.Fatalf("GetLastFive(userB) failed: %v", err)
	}
	if len(gotB) != 1 || gotB[0] != "b-1" {
		t.Fatalf("unexpected userB passwords: %v", gotB)
	}
}
