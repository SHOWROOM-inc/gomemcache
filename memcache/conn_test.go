package memcache

import (
	"testing"
	"time"
)

func TestConn_isExpired(t *testing.T) {
	t.Run("expired by createdAt", func(t *testing.T) {
		now := time.Date(2023, 11, 17, 1, 0, 11, 0, time.UTC)

		cn := &conn{
			createdAt:  time.Date(2023, 11, 17, 1, 0, 0, 0, time.UTC),
			lastUsedAt: now,
			c: &Client{
				ConnMaxLifeTime: 10 * time.Second,
				ConnMaxIdleTime: 5 * time.Second,
			},
		}

		actual := cn.isExpired(now)
		if actual == false {
			t.Fatalf("should be true")
		}
	})
	t.Run("expired by lastUsedAt", func(t *testing.T) {
		now := time.Date(2023, 11, 17, 1, 0, 11, 0, time.UTC)

		cn := &conn{
			createdAt:  time.Date(2023, 11, 17, 1, 0, 2, 0, time.UTC),
			lastUsedAt: time.Date(2023, 11, 17, 1, 0, 5, 0, time.UTC),
			c: &Client{
				ConnMaxLifeTime: 10 * time.Second,
				ConnMaxIdleTime: 5 * time.Second,
			},
		}

		actual := cn.isExpired(now)
		if actual == false {
			t.Fatalf("should be true")
		}
	})
	t.Run("not expired", func(t *testing.T) {
		now := time.Date(2023, 11, 17, 1, 0, 11, 0, time.UTC)

		cn := &conn{
			createdAt:  time.Date(2023, 11, 17, 1, 0, 1, 0, time.UTC),
			lastUsedAt: now,
			c: &Client{
				ConnMaxLifeTime: 10 * time.Second,
				ConnMaxIdleTime: 10 * time.Second,
			},
		}

		actual := cn.isExpired(now)
		if actual == true {
			t.Fatalf("should be false")
		}
	})
}
