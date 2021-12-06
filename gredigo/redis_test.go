package gredigo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestRedisPool(t *testing.T) {
	conf := &RedisConf{
		Host:        "127.0.0.1",
		Port:        6379,
		MaxIdle:     60,
		MaxActive:   200,
		IdleTimeout: 240,
	}

	// 建立连接
	conf.SetRedisPool("default")
	defer func() {
		m := CloseAllPool()
		log.Println("close all pool,result: ", m)
	}()

	var wg sync.WaitGroup

	for i := 0; i < 120; i++ {
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
			_, _ = client.Do("lpush", "my_list", 123)
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
ok
*/

func TestGetRedisClientWithTimeout(t *testing.T) {
	conf := &RedisConf{
		Host:        "127.0.0.1",
		Port:        6379,
		MaxIdle:     100,
		MaxActive:   200,
		IdleTimeout: 240,
	}

	// 建立连接
	conf.SetRedisPool("default")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := GetRedisClientWithTimeout("default", ctx)
	defer client.Close()

	ok, err := client.Do("set", "user_name", "daheige")
	fmt.Println(ok, err)

	value, _ := redis.String(client.Do("get", "user_name"))
	fmt.Println("user_name:", value)

	// 切换到database 1上面操作
	v, err := client.Do("Select", 1)
	fmt.Println(v, err)
	_, _ = client.Do("lpush", "my_list", 123)
	log.Println("exec success...")
}
