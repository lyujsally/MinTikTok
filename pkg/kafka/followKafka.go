package kafka

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lyujsally/MinTikTok-lyujsally/dao/mysql"
	"github.com/lyujsally/MinTikTok-lyujsally/dao/redis"
	"github.com/lyujsally/MinTikTok-lyujsally/settings"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type FollowKafka struct {
	BrokerList     []string
	Topic          string
	FProducer      sarama.SyncProducer
	FConsumerGroup sarama.ConsumerGroup
}

// 创建一个新的FollowKafka实例
func NewFollowKafka(brokerList []string, topic string) *FollowKafka {

	followKafka := &FollowKafka{
		BrokerList: brokerList,
		Topic:      topic,
	}

	// 创建生产者
	followProducer, err := sarama.NewSyncProducerFromClient(KfkCli.conn)
	if err != nil {
		log.Fatalf("Failed to create followProducer: %s", err)
	}

	// 创建消费者组
	//followConsumer, err := sarama.NewConsumerFromClient(KfkCli.conn)
	followConsumerGroup, err := sarama.NewConsumerGroupFromClient(fmt.Sprintf("%s%s", topic, "_group"), KfkCli.conn)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %s", err)
	}

	followKafka.FProducer = followProducer
	followKafka.FConsumerGroup = followConsumerGroup

	return followKafka
}

// 关闭 followkafka 生产者与消费者
func (fk *FollowKafka) CloseFollowKafka() {
	fk.FProducer.Close()
	fk.FConsumerGroup.Close()
}

// follow消息发布
func (fk *FollowKafka) FollowProducer(msg string) error {
	//创建消息
	kafkaMsg := &sarama.ProducerMessage{
		Topic: fk.Topic,
		Value: sarama.StringEncoder(msg),
	}

	//使用生产者发送消息
	partition, offset, err := fk.FProducer.SendMessage(kafkaMsg)
	if err != nil {
		zap.L().Error("Failed to send followMessage")
		return err
	}

	zap.L().Debug("FollowMessage sent successfully",
		zap.Int32("Partition", partition),
		zap.Int64("Offset", offset))

	return nil
}

// follow消息消费
func (fk *FollowKafka) FollowConsumer() error {

	switch fk.Topic {
	case "follow":
		go func() {
			if err := fk.FConsumerGroup.Consume(context.Background(), []string{fk.Topic}, &followConsumerHandler{}); err != nil {
				log.Fatalf("Failed followConsumerHandler: %s", err)
			}
		}()
	case "unfollow":
		go func() {
			if err := fk.FConsumerGroup.Consume(context.Background(), []string{fk.Topic}, &unfollowConsumerHandler{}); err != nil {
				log.Fatalf("Failed unfollowConsumerHandler: %s", err)
			}
		}()
	default:
		log.Fatalf("Invalid topic: %s", fk.Topic)
		return fmt.Errorf("Invalid topic: %s", fk.Topic)
	}

	return nil
}

type followConsumerHandler struct{}

func (h *followConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 初始化操作，可以留空
	log.Printf("follow消费者组成功启动")
	return nil
}

func (h *followConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// 清理操作，可以留空
	return nil
}

func (h *followConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		// 参数解析
		params := strings.Split(fmt.Sprintf("%s", msg.Value), ",")
		userId, _ := strconv.Atoi(params[0])
		targetId, _ := strconv.Atoi(params[1])

		// 日志记录
		sql := fmt.Sprintf("CALL addFollowRelation(%v,%v)", targetId, userId)

		zap.L().Debug("消费队列执行添加关系", zap.String("SQL:", sql))

		//执行SQL
		if err := mysql.DB.Raw(sql).Scan(nil).Error; nil != err {
			// 执行出错，打印日志。
			log.Println(err.Error())
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

type unfollowConsumerHandler struct{}

func (h *unfollowConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 初始化操作，可以留空
	log.Printf("unfollow消费者组成功启动")
	return nil
}

func (h *unfollowConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// 清理操作，可以留空
	return nil
}

func (h *unfollowConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 参数解析
		params := strings.Split(fmt.Sprintf("%s", msg.Value), ",")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		targetId, _ := strconv.ParseInt(params[1], 10, 64)

		log.Printf("%v,%v", userId, targetId)

		// 日志记录
		sql := fmt.Sprintf("CALL delFollowRelation(%v,%v)", targetId, userId)

		zap.L().Debug("消费队列执行删除关系", zap.String("SQL:", sql))

		//执行SQL
		if err := mysql.DB.Raw(sql).Scan(nil).Error; nil != err {
			// 执行出错，打印日志。
			log.Println(err.Error())
		}

		session.MarkMessage(msg, "")

		// 再删Redis里的信息，防止脏数据，保证最终一致性
		redis.ReUpdateRedisUnfollow(userId, targetId)
	}

	return nil
}

var KfkFollowAdd *FollowKafka
var KfkFollowDel *FollowKafka

// 初始化follow相关的Kafka连接
func InitFollowKafka(cfg *settings.KafkaConfig) {

	KfkFollowAdd = NewFollowKafka(cfg.Broker, cfg.Topic2)
	go KfkFollowAdd.FollowConsumer()
	KfkFollowDel = NewFollowKafka(cfg.Broker, cfg.Topic3)
	go KfkFollowDel.FollowConsumer()

}
