package bucket

import (
	"sync"
	"time"

	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
)

type Bucket struct {
	capacity      int
	currentAmount int
	rate          time.Duration
	lastCheck     time.Time
	ticker        *time.Ticker
	doneChannel   chan bool
	mx            sync.RWMutex
	ttl           time.Duration // время жизни бакета без запросов
}

func NewBucket(capacity int, rate time.Duration) (*Bucket, error) {
	if rate == 0 {
		return nil, errors.ErrBucketRateInvalid
	}

	b := &Bucket{
		capacity:      capacity,
		currentAmount: capacity,
		rate:          rate,
		lastCheck:     time.Now(),
		doneChannel:   make(chan bool, 1),
		mx:            sync.RWMutex{},
		ttl:           time.Minute * 2, // default ttl
	}

	go func(b *Bucket) {
		b.ticker = time.NewTicker(b.rate)

		defer b.ticker.Stop()

		for {
			select {
			case <-b.ticker.C:
				b.mx.RLock()
				if time.Since(b.lastCheck) > b.ttl {
					b.doneChannel <- true
					b.mx.RUnlock()

					return
				}
				b.mx.RUnlock()

				if b.currentAmount == b.capacity {
					continue
				}

				b.mx.Lock()
				b.currentAmount++
				b.mx.Unlock()
			default:
				continue
			}
		}
	}(b)

	return b, nil
}

func (b *Bucket) Allow() bool {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.lastCheck = time.Now()

	if b.currentAmount > 0 {
		b.currentAmount--
		return true
	}

	return false
}

func (b *Bucket) Capacity() int {
	return b.capacity
}

func (b *Bucket) Amount() int {
	b.mx.RLock()
	defer b.mx.RUnlock()

	return b.currentAmount
}

func (b *Bucket) Reset() {
	b.mx.Lock()
	b.currentAmount = b.capacity
	b.mx.Unlock()
}

func (b *Bucket) GetDoneChannel() chan bool {
	b.mx.RLock()
	defer b.mx.RUnlock()

	return b.doneChannel
}
