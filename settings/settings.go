package settings

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*KafkaConfig `mapstructure:"kafka"`
	*EtcdConfig  `mapstructure:"etcd"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type KafkaConfig struct {
	Topic1 string   `mapstructure:"topic1"`
	Topic2 string   `mapstructure:"topic2"`
	Topic3 string   `mapstructure:"topic3"`
	Broker []string `mapstructure:"broker"`
}

type EtcdConfig struct {
	ServiceName  string        `mapstructure:"servicename"`
	Endpoints    []string      `mapstructure:"endpoints"`
	SerrviceAddr string        `mapstructure:"serviceaddr"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
}

type ProjectPathConfig struct {
	ProjectPath string `yaml:"project_path"`
}

func Init(fileName string) (err error) {

	viper.SetConfigFile(fileName)

	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		log.Printf("viper.ReadInconfig() failed,err:%v", err)
		return
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("配置文件修改了...")
		if viper.Unmarshal(Conf); err != nil {
			log.Printf("viper.Unmarshal failed,err:%v \n", err)
		}
	})
	return
}
