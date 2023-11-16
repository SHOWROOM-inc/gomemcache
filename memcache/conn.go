package memcache

import (
	"bufio"
	"net"
	"time"
)

// conn is a connection to a server.
type conn struct {
	nc         net.Conn
	rw         *bufio.ReadWriter
	addr       net.Addr
	p          *pool
	c          *Client
	lastUsedAt time.Time
	createdAt  time.Time
}

// release returns this connection back to the client's free pool
func (cn *conn) release() {
	cn.p.putFreeConn(cn)
}

func (cn *conn) close() {
	cn.p.closeConn(cn)
}

func (cn *conn) extendDeadline() {
	cn.nc.SetDeadline(time.Now().Add(cn.c.netTimeout()))
}

func (cn *conn) isExpired(now time.Time) bool {
	if cn.createdAt.Add(cn.c.ConnMaxLifeTime).Before(now) {
		return true
	}
	if cn.lastUsedAt.Add(cn.c.ConnMaxLifeTime).Before(now) {
		return true
	}
	return false
}

// condRelease releases this connection if the error pointed to by err
// is nil (not an error) or is only a protocol level error (e.g. a
// cache miss).  The purpose is to not recycle TCP connections that
// are bad.
func (cn *conn) condRelease(err *error) {
	if *err == nil || resumableError(*err) {
		cn.release()
	} else {
		cn.close()
	}
}
