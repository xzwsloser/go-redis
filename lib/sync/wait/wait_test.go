package wait

import (
	"log"
	"testing"
	"time"
)

func TestWaitGroup(t *testing.T) {
	wt := &Wait{}
	wt.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wt.Done()
			time.Sleep(time.Millisecond * 150)
		}()
	}

	if wt.WaitWithTimeout(time.Second) {
		log.Println("到达超时时间 ...")
	} else {
		log.Println("正常结束")
	}
}
