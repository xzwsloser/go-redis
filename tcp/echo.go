package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoClient struct {
	conn net.Conn
	wait wait.Wait
}

func (ec *EchoClient) Close() error {
	ec.wait.WaitWithTimeout(time.Second * 10)
	return ec.conn.Close()
}

func (eh *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if eh.closing.Get() {
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		conn: conn,
	}
	eh.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Warn("the client is closed ...")
				eh.activeConn.Delete(client)
				client.Close()
			} else {
				logger.Warn("read from client err ", err)
			}

			return
		}

		client.wait.Add(1)
		_, _ = client.conn.Write([]byte(message))
		client.wait.Done()
	}
}

func (eh *EchoHandler) Close() error {
	logger.Info("the error handler is closing ...")
	eh.closing.Set(true)
	eh.activeConn.Range(func(key any, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})

	return nil
}
