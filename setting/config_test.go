package setting

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestViperConfig_Load(t *testing.T) {
	conf := New(WithConfigFile("./test.yaml"), WithWatchFile())
	if err := conf.Load(); err != nil {
		log.Fatalln("load config err: ", err)
	}

	log.Println(conf.IsSet("app"))

	var app int64
	_ = conf.ReadSection("app", &app)
	log.Println(app)

	// check config change
	for {
		_ = conf.ReadSection("app", &app)
		log.Println(app)
		time.Sleep(2 * time.Second)
	}
}

// app.yaml section config.
var (
	appServerConf = &appServerSettingS{}
)

// appServerSettingS server config.
type appServerSettingS struct {
	AppEnv              string        `json:"app_env"`
	AppDebug            bool          `json:"app_debug"`
	GRPCPort            int           `json:"grpc_port"`
	GRPCHttpGatewayPort int           `json:"grpc_http_gateway_port"`
	HttpPort            int           `json:"http_port"`
	ReadTimeout         time.Duration `json:"read_timeout"`
	WriteTimeout        time.Duration `json:"write_timeout"`
	LogDir              string        `json:"log_dir"`
	JobPProfPort        int           `json:"job_p_prof_port"`
}

// readConfig 读取配置文件
func readConfig(configDir string) error {
	// 测试拓展名获取
	filename := "abc.yaml"
	log.Println(strings.TrimPrefix(filepath.Ext(filename), ".")) // yaml

	log.Println(filepath.Dir("/abc/app.yaml"))
	conf := New(WithConfigFile(configDir + "/test.yaml"))
	err := conf.Load()
	if err != nil {
		return err
	}

	err = conf.ReadSection("AppServer", &appServerConf)
	if err != nil {
		return err
	}

	appServerConf.ReadTimeout *= time.Second
	appServerConf.WriteTimeout *= time.Second

	if appServerConf.AppDebug {
		log.Println("app server config: ", appServerConf)
	}

	return nil
}

func TestReadSection(t *testing.T) {
	_ = readConfig("./")
	b, _ := json.Marshal(appServerConf)
	log.Println("section app config: ", string(b))
}

/*
=== RUN   TestReadSection
2023/04/11 17:02:11 yaml
2023/04/11 17:02:11 /abc
2023/04/11 17:02:11 app server config:  &{dev true 50051 1336 1338 6s 6s ./logs 30031}
2023/04/11 17:02:11 section app config:  {"app_env":"dev","app_debug":true,"grpc_port":50051,
"grpc_http_gateway_port":1336,"http_port":1338,"read_timeout":6000000000,"write_timeout":6000000000,
"log_dir":"./logs","job_p_prof_port":30031}
--- PASS: TestReadSection (0.00s)
PASS
*/
