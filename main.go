package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/lyujsally/MinTikTok-lyujsally/dao/mysql"
	"github.com/lyujsally/MinTikTok-lyujsally/dao/redis"
	"github.com/lyujsally/MinTikTok-lyujsally/kitex_gen/relation/relationservice"
	"github.com/lyujsally/MinTikTok-lyujsally/pkg/kafka"

	"github.com/lyujsally/MinTikTok-lyujsally/settings"
)

func main() {

	var filePath string
	flag.StringVar(&filePath, "f", "./conf/config.yaml", "配置文件的路径")
	flag.Parse()
	// 加载配置文件
	err := settings.Init(filePath)
	if err != nil {
		log.Printf("init settings failed,err:%v\n", err)
		return
	}
	// 初始化mysql
	if err = mysql.Init(settings.Conf.MysqlConfig); err != nil {
		log.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 初始化redis
	if err := redis.InitRedisCli(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis client failed, err:%v\n", err)
		return
	}
	fmt.Printf("connect redis success...")
	defer redis.RedisClose()

	//初始化Kafka
	k := kafka.InitKafka(settings.Conf.KafkaConfig)
	defer k.KafkaClose()
	kafka.InitFollowKafka(settings.Conf.KafkaConfig)
	defer kafka.KfkFollowAdd.CloseFollowKafka()

	// 初始化etcd
	r, err := etcd.NewEtcdRegistry(settings.Conf.Endpoints)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", settings.Conf.SerrviceAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	svr := relationservice.NewServer(new(RelationServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: settings.Conf.ServiceName}),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}

}
