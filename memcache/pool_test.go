package memcache

import "testing"

func TestPool_isNewConnOk(t *testing.T) {
	t.Run("MaxOpenConn is zero", func(t *testing.T) {
		p := &pool{
			c: &Client{
				MaxOpenConns: 0,
			},
			openconnsNum: 9999,
		}
		actual := p.isNewConnOk()
		if actual == false {
			t.Fatalf("should be true")
		}
	})

	t.Run("less than max limit", func(t *testing.T) {
		p := &pool{
			c: &Client{
				MaxOpenConns: 100,
			},
			openconnsNum: 99,
		}
		actual := p.isNewConnOk()
		if actual == false {
			t.Fatalf("should be true")
		}
	})

	t.Run("equal max limit", func(t *testing.T) {
		p := &pool{
			c: &Client{
				MaxOpenConns: 100,
			},
			openconnsNum: 100,
		}
		actual := p.isNewConnOk()
		if actual == true {
			t.Fatalf("should be false")
		}
	})
}
