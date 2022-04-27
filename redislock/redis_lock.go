package redislock

import (
	"errors"
	"strings"

	"github.com/go-redis/redis"
)

// DefaultExpire 加锁的key默认过期时间，单位s
var DefaultExpire = 10

var (
	// ErrLockExists redis lock exists
	ErrRedisLockExists = errors.New("redis lock already exists")

	// ErrRedisClientInvalid redis client invalid
	ErrRedisClientInvalid = errors.New("redis client is nil")

	// ErrRedisLockKeyInvalid lock key invalid
	ErrRedisLockKeyInvalid = errors.New("lock key is empty")
)

// lock lock data
type lock struct {
	client *redis.Client // redis连接句柄，支持redis pool连接句柄
	expire int           // 设置加锁key的过期时间
	key    string        // 加锁的key
	val    interface{}   // 加锁的value，可以是int,int64,string
}

// New 实例化redis分布式锁实例对象
func New(client *redis.Client, key string, val interface{}, expire ...int) (*lock, error) {
	if client == nil {
		return nil, ErrRedisClientInvalid
	}

	if key == "" {
		return nil, ErrRedisLockKeyInvalid
	}

	l := &lock{
		key:    strings.Join([]string{"redis_lock", key}, ":"),
		client: client,
		val:    val,
		expire: DefaultExpire,
	}

	if len(expire) > 0 && expire[0] > 0 {
		l.expire = expire[0]
	}

	return l, nil
}

// delScript lua脚本删除一个key保证原子性，采用lua脚本执行
// 保证原子性（redis是单线程），避免del删除了，其他client获得的lock
var delScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`

// Unlock 释放锁采用redis lua脚步执行，成功返回nil
func (l *lock) Unlock() error {
	err := l.client.Eval(delScript, []string{l.key}, l.val).Err()
	return err
}

// TryLock 尝试加锁,如果加锁成功就返回true,nil
// 利用redis set Ex Nx的原子性实现分布式锁
func (l *lock) TryLock() (bool, error) {
	reply, err := l.client.Do("SET", l.key, l.val, "EX", l.expire, "NX").String()
	if err == redis.Nil {
		// The lock was not successful, it already exists.
		return false, ErrRedisLockExists
	}

	if err != nil {
		return false, err
	}

	return reply == "OK", nil
}
