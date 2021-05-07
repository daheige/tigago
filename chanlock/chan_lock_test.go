package chanlock

import (
	"log"
	"sync"
	"testing"
	"time"
)

var count = 1

func TestChanLock(t *testing.T) {
	log.Println("hello")

	var wg sync.WaitGroup

	// 抢占式的更新count，需要对count进行枷锁保护
	// 如果不加锁，count每次执行后，值都不一样
	chLock := NewChanLock()

	nums := 1000
	wg.Add(nums) // 建议一次性实现计数
	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			chLock.Lock()
			defer chLock.Unlock()

			v := count
			log.Println("current count: ", v)
			v++
			count = v
		}()
	}

	log.Println("exec running....")
	wg.Wait()

	log.Println("count: ", count)
}

func TestTryLock(t *testing.T) {
	nums := 100
	var wg sync.WaitGroup
	wg.Add(nums) // 建议一次性实现计数
	chLock := NewChanLock()

	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			if !chLock.TryLock() {
				log.Println("try lock fail")
				return
			}

			defer chLock.Unlock()
			log.Println("lock success")

			time.Sleep(2 * time.Millisecond)
			v := count
			log.Println("current count: ", v)
			v++
			count = v
		}()
	}

	log.Println("exec running....")
	wg.Wait()

	log.Println("count: ", count)
}

/**$ go test -v -test.run TestChanLock
2019/11/27 21:59:50 current count:  1000
2019/11/27 21:59:50 count:  1001
--- PASS: Test_chanLock (0.03s)
PASS
ok
*/
