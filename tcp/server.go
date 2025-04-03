package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// the config of the tcp server
type Config struct {
	Address    string
	MaxConnect uint32
	Timeout    time.Duration
}

var ClientCounter int32

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) {
	sigCh := make(chan os.Signal)
	closeCh := make(chan struct{})
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP:
			closeCh <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Error("failed to listen to address ", cfg.Address)
		return
	}

	ListenAndServe(listener, handler, closeCh)
}

func ListenAndServe(listener net.Listener, handler tcp.Handler,
	closeChan <-chan struct{}) {
	errCh := make(chan error)
	var wg sync.WaitGroup
	go func() {
		select {
		case <-closeChan:
			logger.Info("receive the quit signal")
		case er := <-errCh:
			logger.Error("receive the error ", er.Error())
		}

		_ = listener.Close()
		_ = handler.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			// timeout continue
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				time.Sleep(time.Millisecond * 5)
				continue
			}

			errCh <- err
			break
		}

		atomic.AddInt32(&ClientCounter, 1)
		wg.Add(1)
		go func() {
			defer func() {
				atomic.AddInt32(&ClientCounter, -1)
				wg.Done()
			}()
			handler.Handle(context.Background(), conn)
		}()
	}

	wg.Wait()
}
