package limitedio

import (
	"context"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
	"wget/fileio"
)

// stringReadCloser implements io.ReadCloser by wrapping a string reader
type stringReadCloser struct {
	*strings.Reader
}

// Close implements io.Closer by doing nothing
func (s stringReadCloser) Close() error {
	return nil // No resources to close for a string reader
}

// NewStringReadCloser creates a new io.ReadCloser from a string
func NewStringReadCloser(s string) io.ReadCloser {
	return stringReadCloser{strings.NewReader(s)}
}

// TestLimitedRead aims to assert that when we have a bytes buffer ([]byte), then, we can get an io.Reader
// to read into specific indices in the buffer; say, we have a bytes buffer of 1500 bytes,
// but we only need the io.Reader to read into the first 5 bytes of the buffer
func TestLimitedRead(t *testing.T) {
	// Initialize a reader from strings of 1500 characters (1.5kB)
	reader := strings.NewReader(strings.Repeat("Hello", 300))

	// Create a large enough buffer to read all reader's content
	buffer := make([]byte, 1500)
	// Only read the first 5 bytes onto buffer
	n, err := reader.Read(buffer[:5])
	if err != nil {
		t.Fatal(err)
	}

	// above, we only gave the Reader a slice of 5 bytes, thus, we expect that it only read 5 bytes
	if n != 5 {
		t.Fatalf("expected to read exactly 5 bytes, got %d bytes read", n)
	}

	// Expect that the first 5 bytes in buffer create the string Hello
	if string(buffer[:5]) != "Hello" {
		t.Fatalf("expected to read \"Hello\", got %q", string(buffer[:5]))
	}

	// Expect that all remaining bytes after the first 5 bytes in buffer are all zeroed
	if !reflect.DeepEqual(buffer[5:], make([]byte, 1500-5)) {
		t.Fatalf("expected that all remaining bytes after the first 5 bytes in buffer are all zeroed")
	}
}

func TestNewSGReader(t *testing.T) {
	// Initialize a reader from strings of 1500 characters (1.5kB)
	reader := NewStringReadCloser(strings.Repeat("Hello", 300))

	type args struct {
		speed  int32
		reader *io.ReadCloser
	}
	tests := []struct {
		name      string
		args      args
		want      *SGReader
		willPanic bool
	}{
		{
			name: "Speed limit 0 bytes/second",
			args: args{
				speed:  0,
				reader: &reader,
			},
			want: &SGReader{
				speed:  math.MaxInt32,
				reader: &reader,
			},
			willPanic: false,
		},

		{
			name: "Speed limit 0 bytes/second nil reader",
			args: args{
				speed:  0,
				reader: nil,
			},
			want: &SGReader{
				speed:  math.MaxInt32,
				reader: nil,
			},
			willPanic: true,
		},

		{
			name: "Speed limit -10 bytes/second",
			args: args{
				speed:  -10,
				reader: &reader,
			},
			want: &SGReader{
				speed:  math.MaxInt32,
				reader: &reader,
			},
			willPanic: false,
		},

		{
			name: "Speed limit -10 bytes/second nil reader",
			args: args{
				speed:  -10,
				reader: nil,
			},
			want: &SGReader{
				speed:  math.MaxInt32,
				reader: nil,
			},
			willPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var got *SGReader
				panicked, _ := checkPanic(
					func() {
						got = NewSGReader(tt.args.speed, tt.args.reader)
					},
				)

				if tt.willPanic {
					if !panicked {
						t.Errorf("NewSGReader() did not panic when it should")
					}
				} else {
					if panicked {
						t.Errorf("NewSGReader() panicked when it should not have")
						return
					}

					if got.speed != tt.want.speed {
						t.Errorf("expected speed %d, got %d", tt.want.speed, got.speed)
					}

					if got.reader != tt.want.reader {
						t.Errorf("expected reader %v, got %v", tt.want.reader, got.reader)
					}
				}
			},
		)
	}
}

// TestSGReader_Read tests whether the SGReader limits reads on a string based io.Reader.
// For this test, we set the max bandwidth to 1 byte/second, in which case we
// expect that 5 bytes will be read at no earlier than 5 seconds, and that after
// each second, the reader should have ideally read 1 more byte into our database buffer
func TestSGReader_Read(t *testing.T) {
	s := strings.Repeat("Hello", 1)
	reader := NewStringReadCloser(s)

	// allow reads of strictly upto 1 byte/second
	sgReader := NewSGReader(1, &reader)
	defer fileio.Close(sgReader)

	sgReader.SetRateListener(
		func(speed int32) {
			fmt.Printf("[0] Rate: %d bytes/second\n", speed)
			if speed > 1 {
				t.Errorf("Reader read faster than 1 byte/second -> at %v bytes/second\n", speed)
			}
		},
	)

	readStart := make(chan struct{})

	var database []byte
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := 0
	go func() {
		// Create a ticker that ticks every second
		ticker := time.NewTicker(1 * time.Second)

		//// Cancel the context after 10 seconds
		//go func() {
		//	time.Sleep(10 * time.Second)
		//	cancel()
		//}()

		// Wait for the first read
		fmt.Println("[1] Wait for the first read >>")
		<-readStart
		fmt.Println("[1]  Done.")
		// Loop until the context is canceled
		for {
			select {
			case <-ticker.C:
				// Check that the database has exactly n, n-1, or n+1 bytes, (graceful test)
				l := len(database)
				if !(l == n || l == n-1 || l == n+1) {
					t.Error(
						"[1] Database doesn't have the expected number of bytes; reads exceeding speed limit, " +
							"or slower than speed limit??",
					)
					fmt.Println("[1] Ticker stopped")
					return
				}
				n++
			case <-ctx.Done():
				fmt.Println("[1] Ticker stopped by ctx")
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)

	// Oversize buffer to check edge cases
	buffer := make([]byte, 1024)
	fmt.Println("[0] Sending start >>")
	readStart <- struct{}{}
	fmt.Println("[0]  Done.")
	startTime := time.Now()
	for {
		fmt.Println("[0] Reading >>")
		n, err := sgReader.Read(buffer)
		fmt.Println("[0]  Done.")
		if err != nil {
			fmt.Println("[0] Finished reading", err)
			break
		}

		database = append(database, buffer[:n]...)
	}
	endTime := time.Now()

	cancel()
	if string(database) != s {
		t.Errorf("[0] expected database to be an exact copy of s -> \n%s, got \n%s", s, string(database))
	}

	// we expect that the 5 byte string (at a read rate of 1 byte/second)
	//should have been read at no earlier than 5 seconds,
	//but ideally no more than 6 seconds in the worst case
	duration := endTime.Sub(startTime)
	if duration < 5*time.Second {
		t.Errorf("Reader read faster, took less than 5 seconds -> at %v second(s)", duration)
	} else if duration > 6*time.Second {
		t.Errorf("Reader read too slow, took extremely more than 5 seconds -> at %v second(s)", duration)
	}
}

// checkPanic returns true if the given function, when called, caused a panic
// call, along with the panic message, otherwise returns false
func checkPanic(f func()) (panicked bool, message interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			message = r
		}
	}()
	f()
	return
}
