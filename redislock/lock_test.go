package redislock

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func getClient() *redis.Client {
	opt := &redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		DB:           0, // use default DB
		MaxRetries:   1,
		DialTimeout:  5 * time.Second, // Default is 5 seconds.
		ReadTimeout:  3 * time.Second, // Default is 3 seconds.
		WriteTimeout: 3 * time.Second, // Default is ReadTimeout.
		PoolSize:     10,
		PoolTimeout:  5 * time.Second,
		MinIdleConns: 30,
		IdleTimeout:  5 * time.Minute,
		MaxConnAge:   1800 * time.Second,
	}

	client := redis.NewClient(opt)
	return client
}

// TestRedisLock 测试枷锁操作
func TestRedisLock(t *testing.T) {
	client := getClient()
	l, err := New(client, "daheige", "hello,world", 20)
	if err != nil {
		t.Fatalf("create redis lock instance err:%v", err)
	}

	if ok, err := l.TryLock(); ok {
		log.Println("lock success")
		for i := 0; i < 10; i++ {
			log.Println("hello,i: ", i)
		}

		l.Unlock()
	} else {
		log.Println("lock fail")
		log.Println("err: ", err)
	}
}

// TestLock 并发操作的尝试枷锁
func TestLock(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)

	client := getClient()
	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			l, err := New(client, "heige", "hello,world", 2)
			if err != nil {
				log.Printf("create redis lock instance err:%v\n", err)
				return
			}

			if ok, err := l.TryLock(); ok {
				log.Println("lock success")
				for i := 0; i < 10; i++ {
					log.Println("hello,i: ", i)
				}

				l.Unlock()
			} else {
				log.Println("lock fail")
				log.Println("err: ", err)
			}
		}(&wg)
	}

	wg.Wait()
	log.Println("ok")
}

/**
=== RUN   TestLock
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 lock success
2022/04/27 22:52:01 hello,i:  0
2022/04/27 22:52:01 hello,i:  1
2022/04/27 22:52:01 hello,i:  2
2022/04/27 22:52:01 hello,i:  3
2022/04/27 22:52:01 hello,i:  4
2022/04/27 22:52:01 hello,i:  5
2022/04/27 22:52:01 hello,i:  6
2022/04/27 22:52:01 hello,i:  7
2022/04/27 22:52:01 hello,i:  8
2022/04/27 22:52:01 hello,i:  9
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 lock fail
2022/04/27 22:52:01 err:  redis lock already exists
2022/04/27 22:52:01 ok
--- PASS: TestLock (0.00s)
PASS
*/
