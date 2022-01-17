package pool

import "time"

type options struct {
	capacity int64 // 协程池大小
	expiryDuration time.Duration // worker 过期时间， 过期清理
}

type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

const (
	defaultExpiryDuration = time.Minute // 默认过期时间
	defaultCapacity = 1000
)

func newDefaultOptions() *options {
	return &options{
		capacity: defaultCapacity,
		expiryDuration: defaultExpiryDuration,
	}
}

// WithCapacity set capacity.
func WithCapacity(capacity int64) Option {
	return option(func(opts *options) {
		opts.capacity = capacity
	})
}

// WithExpiryDuration set expiryDuration.
func WithExpiryDuration(expiryDuration time.Duration) Option {
	return option(func(opts *options) {
		opts.expiryDuration = expiryDuration
	})
}