package main

import (
	"context"
	"flag"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// Only checks the actual ticker
func TestMain_InitialAndFirstTick(t *testing.T) {
	// Run with the short option: -interval=5 (seconds)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-interval=5"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var calls int32
	extractFn = func(ctx context.Context) error {
		n := atomic.AddInt32(&calls, 1)
		if n >= 2 {
			// After initial call (n=1) + first ticker tick (n=2), stop.
			_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
		}
		return nil
	}

	done := make(chan struct{})
	go func() {
		main()
		close(done)
	}()

	// Allow enough time for: startup + initial call + ~5s ticker + call + shutdown
	select {
	case <-done:
		// finished quickly (ok)
	case <-time.After(8 * time.Second):
		t.Fatal("main() did not terminate in time (expected after initial+first tick)")
	}

	got := atomic.LoadInt32(&calls)
	if got < 2 {
		t.Fatalf("expected at least 2 calls (initial + first tick), got %d", got)
	}
}
