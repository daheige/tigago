package gredigo

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestRedisPool(t *testing.T) {
	conf := &RedisConf{
		Host:        "127.0.0.1",
		Port:        6379,
		MaxIdle:     100,
		MaxActive:   200,
		IdleTimeout: 240,
	}

	// 建立连接
	conf.SetRedisPool("default")
	var wg sync.WaitGroup

	for i := 0; i < 20000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := GetRedisClient("default")
			defer client.Close()

			ok, err := client.Do("set", "myname", "daheige")
			fmt.Println(ok, err)

			value, _ := redis.String(client.Do("get", "myname"))
			fmt.Println("myname:", value)

			// 切换到database 1上面操作
			v, err := client.Do("Select", 1)
			fmt.Println(v, err)
			_, _ = client.Do("lpush", "myList", 123)
		}()
	}

	wg.Wait()
	log.Println("exec success...")

}

/*
$ go test -v -test.run TestRedisPool
OK <nil>
OK <nil>
OK <nil>
OK <nil>
OK <nil>
OK <nil>
OK <nil>
OK <nil>
2021/05/01 21:52:21 exec success...
--- PASS: TestRedisPool (1.16s)
PASS
ok  	github.com/go-god/gredigo	1.178s
*/
