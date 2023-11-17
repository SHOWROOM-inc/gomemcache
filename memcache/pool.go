package memcache

import (
	"bufio"
	"net"
	"sync"
	"time"
)

func newPool(addr net.Addr, c *Client) *pool {
	return &pool{
		addr:         addr,
		freeconns:    make(chan *conn),
		freeconnsNum: 0,
		openconnsNum: 0,
		c:            c,
		nowFunc:      time.Now,
	}
}

type pool struct {
	lk           sync.Mutex
	addr         net.Addr
	freeconns    chan *conn
	freeconnsNum int
	openconnsNum int
	c            *Client
	nowFunc      func() time.Time
}

func (p *pool) enqueueNewFreeConn() error {
	nc, err := p.c.dial(p.addr)
	if err != nil {
		return err
	}
	now := p.nowFunc()
	newConn := &conn{
		nc:         nc,
		rw:         bufio.NewReadWriter(bufio.NewReader(nc), bufio.NewWriter(nc)),
		addr:       p.addr,
		p:          p,
		c:          p.c,
		lastUsedAt: now,
		createdAt:  now,
	}

	go func() {
		p.freeconnsNum++
		p.openconnsNum++
		p.freeconns <- newConn
	}()

	return nil
}

func (p *pool) getConn() (*conn, error) {
	p.lk.Lock()
	defer p.lk.Unlock()

	if p.freeconnsNum == 0 && p.isNewConnOk() {
		if err := p.enqueueNewFreeConn(); err != nil {
			return nil, err
		}
	}

	cn := <-p.freeconns
	p.freeconnsNum--
	now := p.nowFunc()

	if cn.isExpired(now) {
		p.closeConn(cn)
		return p.getConn()
	}

	cn.lastUsedAt = now
	return cn, nil
}

func (p *pool) isNewConnOk() bool {
	if p.c.MaxOpenConns <= 0 {
		return true
	}
	return p.openconnsNum < p.c.MaxOpenConns
}

func (p *pool) putFreeConn(cn *conn) {
	go func() {
		p.freeconnsNum++
		p.freeconns <- cn
	}()
}

func (p *pool) closeConn(cn *conn) {
	_ = cn.nc.Close()
	p.lk.Lock()
	defer p.lk.Unlock()
	p.openconnsNum--
}
