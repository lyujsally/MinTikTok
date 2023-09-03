package kafka

import (
	"log"

	"github.com/lyujsally/MinTikTok-lyujsally/settings"

	"github.com/Shopify/sarama"
)

type KafkaCli struct {
	conn sarama.Client
}

var KfkCli *KafkaCli

// 初始化Kafka
func InitKafka(cfg *settings.KafkaConfig) *KafkaCli {

	config := sarama.NewConfig()
	//生产者配置
	config.Producer.RequiredAcks = sarama.WaitForAll //设置生产者发送完数据后需要等待所有的leader和follower都确认才视为发送成功
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//消费者配置
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewClient(cfg.Broker, config)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka brokers: %s", err)
	}

	KfkCli = &KafkaCli{
		conn: conn,
	}

	return KfkCli
}

// 关闭Kafka连接
func (k *KafkaCli) KafkaClose() {
	k.conn.Close()
}
