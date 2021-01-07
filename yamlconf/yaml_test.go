package yamlconf

import (
	"log"
	"testing"
	"time"

	"github.com/daheige/tigago/mysql"

	"github.com/daheige/tigago/gredigo"
)

type dataConf struct {
	RedisConf gredigo.RedisConf
	Ip        []string
}

func TestYaml(t *testing.T) {
	conf := NewConf()
	err := conf.LoadConf("test.yaml")
	log.Println(conf.GetData(), err)

	data := conf.GetData()

	var graceful time.Duration
	conf.Get("GracefulWait", &graceful)
	log.Println("graceful: ", graceful)

	log.Println("RedisCommon: ", data["RedisCommon"])

	// 读取数据到结构体中
	var redisConf = &dataConf{}
	conf.GetStruct("RedisCommon", redisConf)
	log.Println(redisConf)
	log.Println("Ip:", redisConf.Ip)
	log.Println(redisConf.RedisConf.Password == "")

	dbConf := &mysql.DbConf{}
	conf.GetStruct("DbDefault", dbConf)
	log.Println(dbConf)
}

/**
 * TestYaml
$ go test -v
=== RUN   TestYaml
2021/01/07 23:47:04 map[] <nil>
2021/01/07 23:47:04 graceful:  5s
2021/01/07 23:47:04 RedisCommon:  <nil>
2021/01/07 23:47:04 &{{127.0.0.1 6379  0 3 10 0 0 0 120 0} [11.12.1.1 11.12.1.2 11.12.1.3]}
2021/01/07 23:47:04 Ip: [11.12.1.1 11.12.1.2 11.12.1.3]
2021/01/07 23:47:04 true
2021/01/07 23:47:04 &{127.0.0.1 3306 root root test   true 10 100 0 0s 0s 0s true
<nil> false true <nil> {  <nil> false 0 false false false} {false <nil> false <nil> <nil>
false false false false false false 0 map[] <nil> <nil> map[] <nil> <nil>} {0s false 0}}
--- PASS: TestYaml (0.00s)
PASS
*/
