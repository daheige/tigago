package wrapper

// Option wrapper option
type Options struct {
	BufCap       int
	RecoveryFunc func()
}

// Options optional function
type Option func(o *Options)

// WithRecover set recover func
func WithRecover(recoveryFunc func()) Option {
	return func(o *Options) {
		o.RecoveryFunc = recoveryFunc
	}
}

// WithBufCap set buf cap
func WithBufCap(c int) Option {
	return func(o *Options) {
		o.BufCap = c
	}
}
