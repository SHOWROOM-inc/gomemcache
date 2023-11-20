package memcache

type Stats struct {
	Address            string
	MaxIdleConnections int
	MaxOpenConnections int
	OpenConnsNum       int
	IdleConnsNum       int
}
