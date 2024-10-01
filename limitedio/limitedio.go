// Package limitedio creates a wrapper around io.Reader to limit the rate of reads from the reader
package limitedio

import (
	"context"
	"io"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// noCopy may be added to structs which must not be copied after the first use.
// See https://golang.org/issues/8005#issuecomment-190753527 for details.
type noCopy struct{}

// waiters keeps track of whether there is some reader waiting to read from the channel, so we can write to it
type waiters struct {
	m    *sync.Mutex
	some bool
	// bidirectional channel to keep readers waiting for the next read allocation after the maximum reads for the
	//current second have been depleted. The read allocator will replenish the reads every second
	channel chan struct{}
}

// init properly initializes the target waiter, in place, with a new mutex and non-nil channel, overriding any values
func (w *waiters) init() {
	w.m = &sync.Mutex{}
	w.channel = make(chan struct{})
	w.some = false
}

// Receive registers a read waiter, who waits to read from the channel, thus,
// writers can write to the channel; then, blocks until it receives some data from the channel
func (w *waiters) Receive() {
	w.m.Lock()
	w.some = true
	w.m.Unlock()

	<-w.channel
}

// Send checks if there is some reader waiting to read from the channel, before sending some data to the channel
func (w *waiters) Send() {
	w.m.Lock()
	defer w.m.Unlock()
	if w.some {
		w.channel <- struct{}{}
	}
}

// SGReader defines the structure of an io.Reader interface, that limits the rate of allowed reads (
// in bytes per second) on the underlying io.Reader. Always use NewSGReader to properly create instances of this struct
type SGReader struct {
	// warn of unnecessary copy of this struct, instead, use pointers to access an instance
	_ noCopy
	// the maximum rate, in bytes per second, that the wrapped reader should be read from
	speed int32
	// pointer to wrapped closable io.Reader
	reader *io.ReadCloser
	// every second we can read a maximum of speed bytes from the wrapped reader.
	//This keeps track of how many bytes can be read before the current second elapses
	reads atomic.Int64
	waiters
	// has the Read function been called since the creation of this struct instance.
	//This helps us to not start tick timers when we don't need them yet; i.e,
	//it is not ideal to display read speed when we haven't even started any reads yet
	called bool
	// callback function to be called just before the underlying reader is closed
	onClose func()
	// see SetRateListener
	rateListener func(speed int32)
}

// NewSGReader creates a new instance of a Speed Governed reader,
// that will allow reads of upto `rate` bytes per second.
//
// Note:
//
// 1. If the `rate` is <= 0, then there will not be any speed limiting
//
// 2. panics if the supplied reader is nil
func NewSGReader(rate int32, reader *io.ReadCloser) *SGReader {
	if rate <= 0 {
		rate = math.MaxInt32
	}

	if reader == nil {
		panic("reader can't be nil")
	}

	sgr := SGReader{
		speed:  rate,
		reader: reader,
	}
	// start by allowing reads upto rate bytes
	sgr.reads.Store(int64(rate))
	sgr.waiters.init()
	return &sgr
}

// SetRateListener sets a callback to be called every second, for the lifetime of this reader,
// with the current read rate, in bytes/second
func (r *SGReader) SetRateListener(rateListener func(rate int32)) {
	r.rateListener = rateListener
}

// once must be called once for every instance of the SGReader to initialize a timer that ensures the read hard limit
// is respected every second, and that the registered rateListener is properly updates every second
func (r *SGReader) once() {
	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Create a ticker that ticks every second
		ticker := time.NewTicker(1 * time.Second)
		// need to stop the timer,
		//and cancel the go routine context when the underlying reader is just about to be closed
		r.onClose = func() {
			ticker.Stop()
			cancel()
		}

		for {
			select {
			case <-ticker.C:
				// replenish the bytes that can be read every second
				old := r.reads.Swap(int64(r.speed))
				rate := int64(r.speed) - old
				if r.rateListener != nil {
					r.rateListener(int32(rate))
				}
				// tell readers that we have allocated some more bandwidth
				r.waiters.Send()
			case <-ctx.Done():
				return
			}
		}

	}()
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// Note that for this SGReader,
// calls to Read may block as more allocations are awaited suppose the read rate for an instance of a second have been
// exceeded.
//
// This implementation honors, the recommendations by the io.ReadCloser Read interface
func (r *SGReader) Read(p []byte) (n int, err error) {
	if !r.called {
		r.called = true
		// Do something that should be done once, just before the first active read
		r.once()
	}

	// read at most r.reads bytes
	bytesToRead := min(int64(len(p)), r.reads.Load())
	if len(p) != 0 && bytesToRead == 0 {
		// There is no available read allocations, wait for the next round of allocation
		r.waiters.Receive()
		bytesToRead = min(int64(len(p)), r.reads.Load())
	}

	// we are about to read some bytes of length `bytesToRead`,
	//remove this from the current allocation
	r.reads.Add(-bytesToRead)

	// Read the specified bytes
	reader := *r.reader
	n, err = reader.Read(p[:bytesToRead])

	return n, err
}

// Close performs cleanup on the SGReader, then closes the underlying closable io.Reader,
// propagating errors as may occur
func (r *SGReader) Close() error {
	// Call registered onclose listener
	if r.onClose != nil {
		r.onClose()
	}

	// Close the underlying reader too
	reader := *r.reader
	return reader.Close()
}
