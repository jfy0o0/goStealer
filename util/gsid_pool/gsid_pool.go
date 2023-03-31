package gsid_pool

type IDPool struct {
	cap   uint64
	queue chan uint64
}

func New(cap uint64) *IDPool {
	pool := &IDPool{
		cap:   cap,
		queue: make(chan uint64, cap),
	}
	pool.init()
	return pool
}

func (p *IDPool) init() {
	var i = uint64(0)
	for ; i < p.cap; i++ {
		p.queue <- i
	}
}

func (p *IDPool) NewID() uint64 {
	return <-p.queue
}

func (p *IDPool) DeleteID(id uint64) {
	p.queue <- id
}

func (p *IDPool) Size() uint64 {
	return uint64(len(p.queue))
}

func (p *IDPool) Cap() uint64 {
	return p.cap
}
