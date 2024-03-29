package gredigo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisConf redis连接信息
// 假如redigo要实现集群参考： redis-go-cluster
type RedisConf struct {
	Host     string
	Port     int
	Password string
	Database int

	// MaxIdle Maximum number of idle connections in the pool.
	MaxIdle int // 空闲连接的最大数量

	// MaxActive Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int // 最大激活数量

	ConnectTimeout int // 连接超时，单位s
	ReadTimeout    int // 读取超时，单位s
	WriteTimeout   int // 写入超时，单位s

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout int // 空闲连接超时,单位s

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime int // 连接最大生命周期,单位s，默认1800s
}

// RedisPoolList 存放连接池信息
var RedisPoolList = map[string]*redis.Pool{}

var (
	// ErrRedisConnectionNotFound redis connection not found
	ErrRedisConnectionNotFound = errors.New("redis connection not found")

	// ErrRedisConnectionInvalid redis connection invalid
	ErrRedisConnectionInvalid = errors.New("redis connection invalid")
)

// GetRedisClient 通过指定name获取池子中的redis连接句柄
func GetRedisClient(name string, timeoutCtx ...context.Context) redis.Conn {
	var ctx context.Context
	if len(timeoutCtx) > 0 && timeoutCtx[0] != nil {
		ctx = timeoutCtx[0]
	} else {
		ctx = context.Background()
	}

	activeConn, err := getRedisConn(name, ctx)
	if err != nil {
		return ErrorConn{err}
	}

	return activeConn
}

// GetRedisClientWithTimeout return redis.ConnWithTimeout
// you can use DoWithTimeout method.
func GetRedisClientWithTimeout(name string, ctx context.Context) redis.ConnWithTimeout {
	pool, ok := RedisPoolList[name]
	if !ok {
		return ErrorConn{ErrRedisConnectionNotFound}
	}

	activeConn, err := pool.GetContext(ctx)
	if err != nil {
		return ErrorConn{err}
	}

	conn, ok := activeConn.(redis.ConnWithTimeout)
	if ok {
		return conn
	}

	return ErrorConn{ErrRedisConnectionInvalid}
}

func getRedisConn(name string, ctx context.Context) (redis.Conn, error) {
	pool, ok := RedisPoolList[name]
	if !ok {
		return nil, ErrRedisConnectionNotFound
	}

	return pool.GetContext(ctx)
}

// 接口静态检测是否实现了redis.Conn
var _ redis.Conn = (*ErrorConn)(nil)

// ErrorConn error connection impl redis.Conn and redis.ConnWithTimeout
type ErrorConn struct{ err error }

func (ec ErrorConn) Do(string, ...interface{}) (interface{}, error) { return nil, ec.err }
func (ec ErrorConn) DoWithTimeout(time.Duration, string, ...interface{}) (interface{}, error) {
	return nil, ec.err
}
func (ec ErrorConn) Send(string, ...interface{}) error                     { return ec.err }
func (ec ErrorConn) Err() error                                            { return ec.err }
func (ec ErrorConn) Close() error                                          { return nil }
func (ec ErrorConn) Flush() error                                          { return ec.err }
func (ec ErrorConn) Receive() (interface{}, error)                         { return nil, ec.err }
func (ec ErrorConn) ReceiveWithTimeout(time.Duration) (interface{}, error) { return nil, ec.err }

// AddRedisPool 添加新的redis连接池
func AddRedisPool(name string, conf *RedisConf) {
	RedisPoolList[name] = NewRedisPool(conf)
}

// SetRedisPool 设置redis连接池
func (r *RedisConf) SetRedisPool(name string) {
	AddRedisPool(name, r)
}

// ClosePoolByName 通过name释放连接池
func ClosePoolByName(name string) error {
	pool, ok := RedisPoolList[name]
	if !ok {
		return ErrRedisConnectionNotFound
	}

	delete(RedisPoolList, name)

	return pool.Close()
}

// CloseAllPool 释放所有的连接池，返回map[name]error
func CloseAllPool() map[string]error {
	m := make(map[string]error, len(RedisPoolList))
	for name, pool := range RedisPoolList {
		m[name] = pool.Close()
		delete(RedisPoolList, name)
	}

	return m
}

// NewRedisPool 创建redis pool连接池
// If Wait is true and the pool is at the MaxActive limit, then Get() waits
// for a connection to be returned to the pool before returning.
//
// TestOnBorrow is an optional application supplied function for checking
// the health of an idle connection before the connection is used again by
// the application. Argument t is the time that the connection was returned
// to the pool. If the function returns an error, then the connection is
// closed.
func NewRedisPool(conf *RedisConf) *redis.Pool {
	if conf.MaxConnLifetime == 0 {
		conf.MaxConnLifetime = 1800
	}

	if conf.ConnectTimeout == 0 {
		conf.ConnectTimeout = 5
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 3
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 3
	}

	return &redis.Pool{
		Wait:            true, // 等待redis connection放入pool池子中
		MaxIdle:         conf.MaxIdle,
		IdleTimeout:     time.Duration(conf.IdleTimeout) * time.Second,
		MaxActive:       conf.MaxActive,
		MaxConnLifetime: time.Duration(conf.MaxConnLifetime) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port),
				redis.DialReadTimeout(time.Duration(conf.ReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(conf.WriteTimeout)*time.Second),
				redis.DialConnectTimeout(time.Duration(conf.ConnectTimeout)*time.Second),
			)

			if err != nil {
				return nil, err
			}

			if len(conf.Password) != 0 {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			// 选择db
			if conf.Database >= 0 {
				if _, err := c.Do("SELECT", conf.Database); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
	}
}
