# tigago

    Go develop commonly used component library package.

# About package
    
    .
    ├── chanlock            chan实现trylock乐观锁
    ├── crypto              常见的md5,sha1,sha1file,aes/des,ecb,openssl_encrypt实现
    ├── def                 为兼容php其他语言而定义的空数组，空对象
    ├── glog                基于mutex乐观锁实现的每天流动式日志，将日志内容直接落地到文件中
    ├── gnsq                go-nsq基本操作封装
    ├── goredis             基于go-redis/redis封装的redis客户端使用函数（支持cluster集群）
    ├── gpprof              pprof性能分析监控封装
    ├── gqueue              通过指定goroutine个数,实现task queue执行器
    ├── grecover            golang panic/recover捕获堆栈信息实现
    ├── gredigo             基于redigo封装而成的go redis辅助函数，方便快速接入redis操作
    ├── gresty              go http client support get,post,delete,patch,put,head,file method
    ├── gtask               golang task在独立协程中调度实现
    ├── gtime               time相关的一些辅助函数
    ├── gutils              file,num,str字符串相关的一些辅助函数,比如Uuid,HTMLSpecialchars,Uniqid等php函数实现
    ├── gxorm               golang xorm客户端简单封装，方便使用
    ├── jsontime            fix gorm/xorm time.Time json encode/decode bug
    ├── logger              基于zap日志库进行一些必要的优化的日志库
    ├── monitor             基于prometheus二次开发、封装的一些函数，主要用于http/job/grpc服务性能监控
    ├── mutexlock           基于sync.Mutex基础上拓展的乐观锁
    ├── mysql               基于go gorm库封装而成的mysql客户端的一些辅助函数
    ├── mytest              tigago 一些单元测试
    ├── redislock           基于redigo实现的redis+lua分布式锁实现
    ├── runner              runner用于按照顺序，执行程序任务操作，可作为cron作业或定时任务
    ├── sem                 指定数量的空结构体缓存通道，实现信息号实现互斥锁
    ├── setting             通过viper+fsnotify实现配置文件读取，支持配置热更新
    ├── work                利用无缓冲chan创建goroutine池来控制一组task的执行
    ├── workpool            workpool工作池实现，对于百万级并发的一些场景特别适用
    ├── xerrors             自定义错误类型，一般用在api/微服务等业务逻辑中，处理错误
    ├── xsort               基于sort标准库封装的sort操作函数
    └── yamlconf            基于yaml+reflect实现yaml文件的读取，一般用在web/job/rpc应用中

# Usage

    golang1.11+版本，可采用go mod机制管理包,需设置goproxy
    go version >= 1.13
    设置goproxy代理
    vim ~/.bashrc添加如下内容:
    export GOPROXY=https://goproxy.io,direct
    或者
    export GOPROXY=https://goproxy.cn,direct
    或者
    export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

    让bashrc生效
    source ~/.bashrc

# Test unit

    $ go test -v
    997: b75567dc6f88412d55576e4b09127d3f
    998: c3923160f2304849734c0907083f7f65
    999: 8b7a6dce56d346b567c65b3493285831
    --- PASS: TestUuid (0.05s)
        uuid_test.go:13: 测试uuid
    PASS
    ok      github.com/daheige/tigago/mytest       15.841s

# License

    MIT

# Goland Ide

[![jetbrains](jetbrains-variant-2.png "jetbrains")](https://jb.gg/OpenSource)
