package srvdriver_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/testcases"
	"github.com/makasim/flowstatesrv/internal/app"
	"github.com/makasim/flowstatesrv/srvdriver"
)

func TestSuite(t *testing.T) {
	s := testcases.Get(func(t *testing.T) flowstate.Driver {
		l, _ := testcases.NewTestLogger(t)

		t.Cleanup(startSrv(t))

		d, err := srvdriver.New(`http://127.0.0.1:8080`, l)
		t.Cleanup(func() {
			if err := d.Shutdown(context.Background()); err != nil {
				t.Fatalf("driver shutdown: %v", err)
			}
		})

		if err != nil {
			t.Fatalf("srvdriver.New: %v", err)
		}

		return d
	})

	s.SetUpDelayer = false
	//s.DisableGoleak()
	s.Test(t)
}

func startSrv(t *testing.T) func() {
	ctx, cancelCtx := context.WithCancel(context.Background())

	cfg := app.Config{}

	runResCh := make(chan error, 1)
	go func() {
		if err := app.New(cfg).Run(ctx); err != nil {
			runResCh <- err
			log.Printf("ERROR: flowstatesrv: app: run: %v", err)
			return
		}
		runResCh <- nil
	}()

	timeoutT := time.NewTimer(time.Second)
	readyT := time.NewTicker(time.Millisecond * 50)

loop:
	for {
		select {
		case <-timeoutT.C:
			t.Fatalf("app not ready within %s", time.Second)
		case <-readyT.C:
			if err := tcpReady(`127.0.0.1:8080`); err != nil {
				continue loop
			}

			break loop
		}
	}

	return func() {
		cancelCtx()
		select {
		case err := <-runResCh:
			if err != nil {
				t.Fatalf("app shutdown error: %v", err)
			}
			return
		case <-time.After(time.Second):
			t.Fatalf("app shutdown timeout")
		}
	}
}

func tcpReady(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
