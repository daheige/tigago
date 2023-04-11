package wrapper

var _ Wrapper = (*wrapChanImpl)(nil)

// wrapChanImpl wrapper impl
type wrapChanImpl struct {
	bufCap       int
	bufCh        chan struct{}
	recoveryFunc func()
}

// NewChanWrapper create wrapChanImpl entity
// If the wrapper using the chan method needs to specify the number of
// goroutines to be executed,the WithBufCap method needs to be called.
// Otherwise, after the Wait method is executed, some goroutines
// will exit without execution.
func NewChanWrapper(opts ...Option) Wrapper {
	w := &wrapChanImpl{}
	option := &Options{}
	for _, o := range opts {
		o(option)
	}

	if option.BufCap == 0 {
		panic("chan wrapper buf cap must be gt 0")
	}

	w.recoveryFunc = option.RecoveryFunc
	if w.recoveryFunc == nil {
		w.recoveryFunc = defaultRecovery
	}

	w.bufCap = option.BufCap
	w.bufCh = make(chan struct{}, w.bufCap)

	return w
}

// Wrap exec func in goroutine without recover catch
func (c *wrapChanImpl) Wrap(fn func()) {
	go func() {
		defer c.done()
		fn()
	}()
}

// WrapWithRecover safely execute func in goroutine
func (c *wrapChanImpl) WrapWithRecover(fn func()) {
	go func() {
		defer c.recoveryFunc()
		defer c.done()
		fn()
	}()
}

// Wait wait all goroutine finish
func (c *wrapChanImpl) Wait() {
	for i := 0; i < c.bufCap; i++ {
		<-c.bufCh
	}
}

func (c *wrapChanImpl) done() {
	c.bufCh <- struct{}{}
}
