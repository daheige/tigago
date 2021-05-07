package mutexlock

import (
	"log"
	"sync"
	"testing"
)

func TestTryLock(t *testing.T) {
	t.Log("test trylock")
	var mutex = NewMutexLock()
	if mutex.TryLock() {
		t.Log("加锁成功!")
		mutex.Unlock()
	}

	mutex.Lock()
	t.Log("haha")
	mutex.Unlock()

}

/*
TestRace test data race
=== RUN   TestRace
2021/05/07 21:46:31 lock success
2021/05/07 21:46:31 lock fail
2021/05/07 21:46:31 lock fail
2021/05/07 21:46:31 lock success
2021/05/07 21:46:31 lock fail
2021/05/07 21:46:31 lock fail
2021/05/07 21:46:31 lock fail
2021/05/07 21:46:31 x:  52
--- PASS: TestRace (0.00s)
PASS
*/
func TestRace(t *testing.T) {
	var mu Mutex
	var x int
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if mu.TryLock() {
					log.Println("lock success")
					x++
					mu.Unlock()
				} else {
					log.Println("lock fail")
				}
			}()

			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			x++
			mu.Unlock()
		}()
	}

	wg.Wait()

	log.Println("x: ", x)
}

/**
$ go test -v
=== RUN   TestTryLock
--- PASS: TestTryLock (0.00s)
    trylock_test.go:8: test trylock
    trylock_test.go:11: 加锁成功!
    trylock_test.go:16: haha
=== RUN   TestRace
--- PASS: TestRace (0.00s)
PASS
ok      github.com/daheige/tigago/mutexlock    0.003s
*/
