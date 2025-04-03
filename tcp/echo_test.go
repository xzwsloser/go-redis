package tcp

import (
	"bufio"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	closeCh := make(chan struct{})
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Error(err)
		return
	}

	go ListenAndServe(listener, NewEchoHandler(), closeCh)

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 10; i++ {
		valStr := strconv.Itoa(rand.Int()) + "\n"
		_, _ = conn.Write([]byte(valStr))
		reader := bufio.NewReader(conn)
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Error(err)
			return
		}

		if string(line) != valStr {
			t.Error("failed to get echo message")
			return
		}
	}

	_ = conn.Close()
	for i := 0; i < 5; i++ {
		_, _ = net.Dial("tcp", ":8080")
	}

	closeCh <- struct{}{}
	time.Sleep(time.Second)
}
