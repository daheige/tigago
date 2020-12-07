package mysql

import (
	"errors"
	"log"
	"sync"
	"testing"

	"gorm.io/gorm"
)

/**
* sql:
* CREATE DATABASE IF NOT EXISTS test default charset utf8mb4;
* create table user (id int primary key auto_increment,name varchar(200)) engine=innodb;
* 模拟数据插入
* mysql> insert into user (name) values("xiaoming");
   Query OK, 1 row affected (0.11 sec)

   mysql> insert into user (name) values("hello");
   Query OK, 1 row affected (0.04 sec)
*/
type myUser struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(200)"`
}

func (myUser) TableName() string {
	return "user"
}

func TestGorm(t *testing.T) {
	dbConf := &DbConf{
		Ip:           "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "root1234",
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ParseTime:    true,
		ShowSql:      true,
	}

	// 设置db engine name
	err := dbConf.SetDbPool() // 建立db连接池
	log.Println("err: ", err)

	err = dbConf.SetEngineName("default") // 为每个db设置一个engine name
	log.Println("err: ", err)

	// defer dbConf.Close() // 关闭当前连接
	defer CloseAllDb() // 关闭所有的连接句柄

	db, err := GetDbObj("default")
	if err != nil {
		t.Log("get db error: ", err.Error())
		return
	}

	user := &myUser{}
	err = db.Where("name = ?", "hello").First(user).Error
	log.Println("user: ", user, err)

	nums := 100
	var wg sync.WaitGroup
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			db, err := GetDbObj("default")
			if err != nil {
				log.Println("get db error: ", err.Error())
				return
			}

			user := &myUser{}
			db.Where("name = ?", "hello").First(user)
			log.Println(user)
		}()
	}

	wg.Wait()
	log.Println("test success")
}

func TestShortConnect(t *testing.T) {
	getDb := func() (*gorm.DB, error) {
		conf := &DbConf{
			Ip:        "127.0.0.1",
			Port:      3306,
			User:      "root",
			Password:  "root1234",
			Database:  "test",
			ParseTime: true,
			ShowSql:   true,
		}

		// 连接gorm.DB实例对象，并非立即连接db,用的时候才会真正的建立连接
		err := conf.ShortConnect()
		if err != nil {
			return nil, errors.New("set gorm.DB failed")
		}

		return conf.Db(), nil
	}

	// 这里我设置了db max_connections最大连接为1000
	var wg sync.WaitGroup
	// var maxConnections = 30
	var maxConnections = 1000
	// var maxConnections = 2000
	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			db, err := getDb()
			if err != nil {
				log.Println("get db error: ", err.Error())
				return
			}

			// 关闭数据句柄
			defer func() {
				sqlDB, err := db.DB()
				if err != nil {
					log.Println("get db instance error: ", err.Error())
					return
				}

				sqlDB.Close()
			}()

			user := &myUser{}
			db.Where("name = ?", "hello").First(user)
			log.Println(user)
		}()
	}

	wg.Wait()
	log.Println("test success")
}

/** go test -v -test.run TestGorm
采用长连接测试
2019/07/20 20:14:51 &{2 hello}
2019/07/20 20:14:51 &{2 hello}
2019/07/20 20:14:51 test success
--- PASS: TestGorm (0.38s)
PASS
ok      github.com/daheige/thinkgo/mysql        0.388s

采用短连接方式测试
$ go test -v -test.run TestShortConnect
2019/07/20 20:16:23 test success
--- PASS: TestShortConnect (1.22s)
PASS
ok      github.com/daheige/thinkgo/mysql        1.229s

当我们把maxConnections 调到2000后
$ go test -v -test.run TestShortConnect
=== RUN   TestShortConnect
2019/07/20 12:17:51 get db error:  set gorm.DB failed
2019/07/20 12:17:51 get db error:  set gorm.DB failed
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x88 pc=0x6be466]

goroutine 1401 [running]:
就会出现panic

// 对github.com/jinzhu/gorm升级到gorm.io/gorm v1.20.1版本后，单元测试
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 &{5 hello}
2020/09/12 14:02:29 test success
--- PASS: TestGorm (0.03s)
PASS
*/
